/*
Copyright 2019 Baidu, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package edgehandler

import (
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"k8s.io/klog"

	otev1 "github.com/baidu/ote-stack/pkg/apis/ote/v1"
	"github.com/baidu/ote-stack/pkg/clustermessage"
	"github.com/baidu/ote-stack/pkg/clustershim"
	"github.com/baidu/ote-stack/pkg/config"
	oteclient "github.com/baidu/ote-stack/pkg/generated/clientset/versioned"
	"github.com/baidu/ote-stack/pkg/tunnel"
)

var (
	edgeTunnelMsg = []byte("msg")
	LastSend      clustermessage.ClusterMessage
)

type fakeEdgeTunnel struct {
}

type fakeShimHandler struct {
}

func (f *fakeEdgeTunnel) Send(data []byte) error {
	msg := &clustermessage.ClusterMessage{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return err
	}
	LastSend = *msg
	return nil
}

func (f *fakeEdgeTunnel) RegistReceiveMessageHandler(tunnel.TunnelReadMessageFunc) {
	return
}

func (f *fakeEdgeTunnel) RegistAfterConnectToHook(tunnel.AfterConnectToHook) {
	return
}

func (f *fakeEdgeTunnel) Start() error {
	return nil
}

func (f *fakeEdgeTunnel) Stop() error {
	return nil
}

func (f *fakeShimHandler) Do(in *clustermessage.ClusterMessage) (*clustermessage.ClusterMessage, error) {
	head := &clustermessage.MessageHead{
		MessageID:           in.Head.MessageID,
		Command:             clustermessage.CommandType_ControlResp,
		ParentClusterName:   in.Head.ParentClusterName,      
	}

	resp := &clustermessage.ControllerTaskResponse{
            Timestamp:  time.Now().Unix(),
            StatusCode: 200,
			Body:       []byte(""),
	}

	data, err := proto.Marshal(resp)		
    if err != nil {
		klog.Errorf("shim resp to controller task resp failed: %v", err)
		return &clustermessage.ClusterMessage{Head: head,}, nil
	}
	
	msg := &clustermessage.ClusterMessage{
		Head: head,
		Body: data,
	}

	return msg, nil
}

func newFakeShim() clustershim.ShimServiceClient {
	handlers := clustershim.ShimHandler{}
	handlers[otev1.ClusterControllerDestAPI] = &fakeShimHandler{}
	return clustershim.NewlocalShimClientWithHandler(handlers)
}

func TestValid(t *testing.T) {
	succescase := []struct {
		Name string
		Conf *config.ClusterControllerConfig
	}{
		{
			Name: "edgehandler with k8sclient",
			Conf: &config.ClusterControllerConfig{
				ClusterName:           "child",
				ClusterUserDefineName: "child",
				K8sClient:             &oteclient.Clientset{},
				RemoteShimAddr:        "",
				ParentCluster:         "127.0.0.1:8287",
			},
		},
		{
			Name: "edgehandler with remoteshim",
			Conf: &config.ClusterControllerConfig{
				ClusterName:           "child",
				ClusterUserDefineName: "child",
				K8sClient:             nil,
				RemoteShimAddr:        ":8262",
				ParentCluster:         "127.0.0.1:8287",
			},
		},
	}

	for _, sc := range succescase {
		edge := edgeHandler{conf: sc.Conf}
		if err := edge.valid(); err != nil {
			t.Errorf("[%q] unexpected error %v", sc.Name, err)
		}
	}

	errorcase := []struct {
		Name string
		Conf *config.ClusterControllerConfig
	}{
		{
			Name: "cluster name not set",
			Conf: &config.ClusterControllerConfig{
				ClusterName:    "",
				K8sClient:      nil,
				RemoteShimAddr: ":8262",
				ParentCluster:  "127.0.0.1:8287",
			},
		},
		{
			Name: "shim address not set",
			Conf: &config.ClusterControllerConfig{
				ClusterName:    "child1",
				K8sClient:      nil,
				RemoteShimAddr: "",
				ParentCluster:  "127.0.0.1:8287",
			},
		},
		{
			Name: "ParentCluster not set",
			Conf: &config.ClusterControllerConfig{
				ClusterName:    "child1",
				K8sClient:      nil,
				RemoteShimAddr: ":8262",
				ParentCluster:  "",
			},
		},
	}

	for _, ec := range errorcase {
		edge := &edgeHandler{conf: ec.Conf}
		if err := edge.valid(); err == nil {
			t.Errorf("[%q] expected error", ec.Name)
		}
	}
}

func TestIsRemoteShim(t *testing.T) {
	casetest := []struct {
		Name   string
		Conf   *config.ClusterControllerConfig
		Expect bool
	}{
		{
			Name: "use remote shim",
			Conf: &config.ClusterControllerConfig{
				ClusterName:    "child",
				RemoteShimAddr: ":8262",
				K8sClient:      &oteclient.Clientset{},
			},
			Expect: true,
		},
		{
			Name: "use local shim",
			Conf: &config.ClusterControllerConfig{
				ClusterName:    "child",
				RemoteShimAddr: "",
				K8sClient:      &oteclient.Clientset{},
			},
			Expect: false,
		},
	}
	for _, ct := range casetest {
		edge := &edgeHandler{
			conf: ct.Conf,
		}
		res := edge.isRemoteShim()
		if res != ct.Expect {
			t.Errorf("[%q] expected %v, got %v", ct.Name, ct.Expect, res)
		}
	}
}

func TestSendMessageToTunnel(t *testing.T) {
	conf := &config.ClusterControllerConfig{
		ClusterName:       "child",
		K8sClient:         nil,
		RemoteShimAddr:    ":8262",
		ParentCluster:     "127.0.0.1:8287",
		ClusterToEdgeChan: make(chan clustermessage.ClusterMessage),
	}

	controllerAPITask := &clustermessage.ControllerTask{
		Destination: "api",
	}
	controllerAPITaskData, err := proto.Marshal(controllerAPITask)
	assert.Nil(t, err)
	assert.NotNil(t, controllerAPITaskData)

	casetest := []struct {
		Name     string
		SendData clustermessage.ClusterMessage
	}{
		{
			Name: "valid send clusterController",
			SendData: clustermessage.ClusterMessage{
				Head: &clustermessage.MessageHead{
					ParentClusterName: "root",
					ClusterSelector:   "c1,c2",
					Command:           clustermessage.CommandType_ControlReq,
				},
				Body: controllerAPITaskData,
			},
		},
	}

	for _, ct := range casetest {
		edge := &edgeHandler{
			conf:       conf,
			edgeTunnel: &fakeEdgeTunnel{},
		}
		go edge.sendMessageToTunnel()
		edge.conf.ClusterToEdgeChan <- ct.SendData
		time.Sleep(1 * time.Second)
		assert.True(t, proto.Equal(&ct.SendData, &LastSend))
	}
}

func TestReceiveMessageFromTunnel(t *testing.T) {
	conf := &config.ClusterControllerConfig{
		ClusterName:       "child",
		K8sClient:         nil,
		RemoteShimAddr:    ":8262",
		ParentCluster:     "127.0.0.1:8287",
		EdgeToClusterChan: make(chan clustermessage.ClusterMessage, 10),
	}

	edge := &edgeHandler{
		conf:       conf,
		edgeTunnel: &fakeEdgeTunnel{},
		shimClient: newFakeShim(),
	}

	controllerAPITask := &clustermessage.ControllerTask{
		Destination: "api",
	}
	controllerAPITaskData, err := proto.Marshal(controllerAPITask)
	assert.Nil(t, err)
	assert.NotNil(t, controllerAPITaskData)

	casetest := []struct {
		Name         string
		Data         *clustermessage.ClusterMessage
		ExpectHandle bool
	}{
		{
			Name: "match rule",
			Data: &clustermessage.ClusterMessage{
				Head: &clustermessage.MessageHead{
					ParentClusterName: "root",
					ClusterSelector:   "c1,c2,child",
					Command:           clustermessage.CommandType_ControlReq,
				},
				Body: controllerAPITaskData,
			},
			ExpectHandle: true,
		},
		{
			Name: "not match rule",
			Data: &clustermessage.ClusterMessage{
				Head: &clustermessage.MessageHead{
					ParentClusterName: "root",
					ClusterSelector:   "c1,c2",
					Command:           clustermessage.CommandType_ControlReq,
				},
				Body: controllerAPITaskData,
			},
			ExpectHandle: false,
		},
	}

	for _, ct := range casetest {
		LastSend.Head.Command = clustermessage.CommandType_Reserved
		msg, err := proto.Marshal(ct.Data)
		assert.Nil(t, err)
		edge.receiveMessageFromTunnel(conf.ClusterName, msg)

		var broadcast clustermessage.ClusterMessage
		go func() {
			broadcast = <-edge.conf.EdgeToClusterChan
		}()

		time.Sleep(1 * time.Second)

		ok := LastSend.Head.Command == clustermessage.CommandType_ControlResp
		assert.Equal(t, ct.ExpectHandle, ok)
		assert.True(t, proto.Equal(ct.Data, &broadcast))
	}
}

func TestHandleMessage(t *testing.T) {
	conf := &config.ClusterControllerConfig{
		ClusterName:       "child",
		K8sClient:         nil,
		RemoteShimAddr:    ":8262",
		ParentCluster:     "127.0.0.1:8287",
		EdgeToClusterChan: make(chan clustermessage.ClusterMessage, 10),
	}
	edge := &edgeHandler{
		conf:       conf,
		edgeTunnel: &fakeEdgeTunnel{},
		shimClient: newFakeShim(),
	}

	casetest := []struct {
		Name         string
		Data         clustermessage.ClusterMessage
		ExpectHandle bool
	}{
		{
			Name: "dispatch to route",
			Data: clustermessage.ClusterMessage{
				Head: &clustermessage.MessageHead{
					ParentClusterName: "root",
					Command:           clustermessage.CommandType_NeighborRoute,
				},
			},
			ExpectHandle: false,
		},
		{
			Name: "dispatch to api",
			Data: clustermessage.ClusterMessage{
				Head: &clustermessage.MessageHead{
					ParentClusterName: "root",
					Command:           clustermessage.CommandType_ControlReq,
				},
			},
			ExpectHandle: true,
		},
		{
			Name: "dispatch to api",
			Data: clustermessage.ClusterMessage{
				Head: &clustermessage.MessageHead{
					ParentClusterName: "root",
					Command:           clustermessage.CommandType_ControlReq,
				},
			},
			ExpectHandle: true,
		},
	}

	for _, ct := range casetest {
		LastSend.Head.Command = clustermessage.CommandType_Reserved
		if err := edge.handleMessage(&ct.Data); err != nil {
			t.Errorf("[%q] unexpected error %v", ct.Name, err)
		}

		time.Sleep(2 * time.Second)
		ok := LastSend.Head.Command == clustermessage.CommandType_ControlResp
		assert.Equal(t, ct.ExpectHandle, ok)
	}

}
