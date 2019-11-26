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

package clustershim

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/golang/protobuf/proto"

	otev1 "github.com/baidu/ote-stack/pkg/apis/ote/v1"
	"github.com/baidu/ote-stack/pkg/clustermessage"
	"github.com/baidu/ote-stack/pkg/config"
	oteclient "github.com/baidu/ote-stack/pkg/generated/clientset/versioned/fake"
	"github.com/baidu/ote-stack/pkg/clustershim/handler"
)

func fakeNewlocalShimClient(c *config.ClusterControllerConfig) ShimServiceClient {
	local := &localShimClient{
		handlers: make(map[string]handler.Handler),
	}
	local.handlers["api"] = &fakeShimHandler{}
	local.handlers[otev1.ClusterControllerDestHelm] = handler.NewHTTPProxyHandler(c.HelmTillerAddr)
	return local
}

func getControllerTask(des string, method string, uri string, t *testing.T) []byte {
	msg := &clustermessage.ControllerTask{
		Destination: des,
		Method:      method,
		URI:         uri,
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Errorf("to controller task request failed: %v", err)
		return nil
	}
	return data
}

func TestShimClientDoControlRequest(t *testing.T) {
	c := &config.ClusterControllerConfig{
		K8sClient: oteclient.NewSimpleClientset(
			&otev1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "c1",
				},
			},
		),
		HelmTillerAddr: "",
	}
	localClient := fakeNewlocalShimClient(c).(*localShimClient)
	assert.Nil(t, localClient.ReturnChan())

	//supportable handler
	des := "api"
	method := "GET"
	uri := "/apis/ote.baidu.com/v1/namespaces/default/clusters"
	data1 := getControllerTask(des, method, uri, t)
	msg1 := clustermessage.ClusterMessage{
		Head: &clustermessage.MessageHead{
			ParentClusterName:     "c1",
			Command:               clustermessage.CommandType_ControlReq,
			ClusterName:           "c1",
		},
		Body: data1,
	}
	resp, err := localClient.DoControlRequest(&msg1)
	assert.Equal(t, clustermessage.CommandType_ControlResp, resp.Head.Command) // local shim client return not nil resp

	// unsupportable handler
	des = "nohandler"
	data2 := getControllerTask(des, method, uri, t)
	
	msg2 := clustermessage.ClusterMessage{
		Head: &clustermessage.MessageHead{
			ParentClusterName:     "c1",
			Command:               clustermessage.CommandType_ControlReq,
			ClusterName:           "c1",
		},
		Body: data2,
	}
	resp, err = localClient.DoControlRequest(&msg2)
	assert.NotNil(t, resp) // local shim client return not nil resp
	assert.NotNil(t, err)

	//unsupportable command
	des = "api"
	data3 := getControllerTask(des, method, uri, t)
	msg3 := clustermessage.ClusterMessage{
		Head: &clustermessage.MessageHead{
			ParentClusterName:     "c1",
			Command:               clustermessage.CommandType_NeighborRoute,
			ClusterName:           "c1",
		},
		Body: data3,
	}
	resp, err = localClient.DoControlRequest(&msg3)
	assert.Nil(t, resp) // local shim client return nil resp
	assert.NotNil(t, err)
}

func TestShimClientDo (t *testing.T) {
	c := &config.ClusterControllerConfig{
		K8sClient: oteclient.NewSimpleClientset(
			&otev1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "c1",
				},
			},
		),
		HelmTillerAddr: "",
	}
	localClient := NewlocalShimClient(c)
	assert.Nil(t, localClient.ReturnChan())
	
	msg := clustermessage.ClusterMessage{
		Head: &clustermessage.MessageHead{
			Command:	clustermessage.CommandType_NeighborRoute,
		},
	}

	resp, err := localClient.Do(&msg)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
}