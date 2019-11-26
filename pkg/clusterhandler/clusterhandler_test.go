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

package clusterhandler

import (
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	otev1 "github.com/baidu/ote-stack/pkg/apis/ote/v1"
	"github.com/baidu/ote-stack/pkg/clustermessage"
	"github.com/baidu/ote-stack/pkg/clusterrouter"
	"github.com/baidu/ote-stack/pkg/config"
	oteclient "github.com/baidu/ote-stack/pkg/generated/clientset/versioned/fake"
	"github.com/baidu/ote-stack/pkg/tunnel"
)

var (
	fakeTunn = newFakeCloudTunnel()
)

func TestInit(t *testing.T) {
	c := &config.ClusterControllerConfig{}
	client := oteclient.NewSimpleClientset()

	h := &clusterHandler{
		conf:      c,
		k8sEnable: false,
	}
	h.conf.ClusterUserDefineName = ""
	assert.Error(t, h.valid())

	// test root cluster config
	h.conf.ClusterUserDefineName = config.RootClusterName
	assert.Error(t, h.valid())
	h.conf.TunnelListenAddr = ":8272"
	assert.Error(t, h.valid())
	h.conf.K8sClient = client
	assert.False(t, h.k8sEnable)
	assert.NoError(t, h.valid())
	assert.True(t, h.k8sEnable)
	h.conf.ParentCluster = "parent"
	assert.Error(t, h.valid())

	// test no-root cluster config
	h.conf.ClusterUserDefineName = "c1"
	h.conf.ParentCluster = ""
	assert.Error(t, h.valid())
	h.conf.ParentCluster = "parent"
	assert.NoError(t, h.valid())

	h2, err := NewClusterHandler(h.conf)
	assert.Nil(t, err)
	assert.NotNil(t, h2)
}

/*
TestSelectChild test clusterHandler.selectChild func.
route used as following:
root
- c1
	- c2
	- c3
- c4
	- c5
*/
func TestSelectChild(t *testing.T) {
	clusterrouter.Router().AddRoute("c1", "c1")
	clusterrouter.Router().AddRoute("c2", "c1")
	clusterrouter.Router().AddRoute("c3", "c1")
	clusterrouter.Router().AddRoute("c4", "c4")
	clusterrouter.Router().AddRoute("c5", "c4")

	msg := &clustermessage.ClusterMessage{
		Head: &clustermessage.MessageHead{
			ClusterSelector: "c3,c5",
		},
	}
	selected := selectChild(msg)
	assert.Equal(t, 2, len(selected))
	assert.Equal(t, "c3", selected["c1"].Head.ClusterSelector)
	assert.Equal(t, "c5", selected["c4"].Head.ClusterSelector)
}

func TestHasToProcessClusterController(t *testing.T) {
	now := time.Now().Unix()
	cc := &otev1.ClusterController{
		ObjectMeta: metav1.ObjectMeta{
			CreationTimestamp: metav1.NewTime(time.Unix(now-1*60*60, 0)),
		},
	}
	assert.False(t, hasToProcessClusterController(cc))

	cc.ObjectMeta.CreationTimestamp = metav1.NewTime(time.Unix(now, 0))
	assert.True(t, hasToProcessClusterController(cc))

	cc.Status = map[string]otev1.ClusterControllerStatus{
		"c1": {},
	}
	assert.False(t, hasToProcessClusterController(cc))
}

func TestSendToChild(t *testing.T) {
	c := &clusterHandler{}
	fakeTunn.reset()
	c.tunn = fakeTunn
	c.sendToChild(nil)
	assert.False(t, fakeTunn.broadcastCalled)
	assert.False(t, fakeTunn.sendCalled)
	c.sendToChild(&clustermessage.ClusterMessage{})
	time.Sleep(1 * time.Second)
	assert.True(t, fakeTunn.broadcastCalled)
	assert.False(t, fakeTunn.sendCalled)
	fakeTunn.reset()
	c.sendToChild(&clustermessage.ClusterMessage{}, "")
	time.Sleep(1 * time.Second)
	assert.False(t, fakeTunn.broadcastCalled)
	assert.True(t, fakeTunn.sendCalled)
}

//func TestAddClusterController(t *testing.T) {
//	c := &clusterHandler{
//		tunn: newFakeCloudTunnel(),
//	}
//	cc := &otev1.ClusterController{
//		ObjectMeta: metav1.ObjectMeta{
//			CreationTimestamp: metav1.NewTime(time.Now()),
//		},
//	}
//	c.addClusterController()
//}

func TestStart(t *testing.T) {
	fakeK8sClient := oteclient.NewSimpleClientset()
	fakeTunn.reset()
	c := &clusterHandler{
		conf: &config.ClusterControllerConfig{
			TunnelListenAddr:      "fake",
			K8sClient:             fakeK8sClient,
			EdgeToClusterChan:     make(chan clustermessage.ClusterMessage),
			ClusterToEdgeChan:     make(chan clustermessage.ClusterMessage),
			ClusterUserDefineName: config.RootClusterName,
		},
		tunn:      fakeTunn,
		k8sEnable: false,
	}
	err := c.valid()
	assert.Nil(t, err)
	err = c.Start()
	assert.Nil(t, err)

	// test msg from parent
	msg := clustermessage.ClusterMessage{
		Head: &clustermessage.MessageHead{
			ClusterName:     "cc1",
			ClusterSelector: "c1",
		},
	}
	clusterrouter.Router().AddRoute("c1", "c1")
	c.conf.EdgeToClusterChan <- msg
	time.Sleep(1 * time.Second)
	assert.False(t, fakeTunn.broadcastCalled)
	assert.True(t, fakeTunn.sendCalled)
}

func TestAfterClusterConnect(t *testing.T) {
	c := &clusterHandler{
		tunn: fakeTunn,
	}
	cr := &config.ClusterRegistry{
		Name: "c1",
	}
	clusterrouter.Router().DelChild("c1", c.sendToChild)
	fakeTunn.reset()
	c.afterClusterConnect(cr)
	time.Sleep(1 * time.Second)
	assert.True(t, fakeTunn.broadcastCalled)
	assert.False(t, fakeTunn.sendCalled)
}

func TestHandleMessageFromChild(t *testing.T) {
	c := newFakeRootClusterHandler(t)

	err := c.handleMessageFromChild("c1", []byte("hahaha"))
	assert.NotNil(t, err)

	msg := &clustermessage.ClusterMessage{
		Head: &clustermessage.MessageHead{},
	}
	// regist cluster without cluster info
	msg.Head.Command = clustermessage.CommandType_ClusterRegist
	ccbytes, err := proto.Marshal(msg)
	assert.Nil(t, err)
	err = c.handleMessageFromChild("c1", ccbytes)
	assert.NotNil(t, err)
	// unregist cluster without cluster info
	msg.Head.Command = clustermessage.CommandType_ClusterUnregist
	ccbytes, err = proto.Marshal(msg)
	assert.Nil(t, err)
	err = c.handleMessageFromChild("c1", ccbytes)
	assert.NotNil(t, err)
	// resp from child with no namespace and name
	msg.Head.ParentClusterName = c.conf.ClusterUserDefineName
	msg.Head.Command = clustermessage.CommandType_ControlResp
	ccbytes, err = proto.Marshal(msg)
	assert.Nil(t, err)
	err = c.handleMessageFromChild("c1", ccbytes)
	assert.Nil(t, err)
	// resp from child transmit to parent
	msg.Head.ParentClusterName = c.conf.ClusterUserDefineName + "1"
	msg.Head.Command = clustermessage.CommandType_ControlResp
	ccbytes, err = proto.Marshal(msg)
	assert.Nil(t, err)
	err = c.handleMessageFromChild("c1", ccbytes)
	assert.Nil(t, err)
}

func TestHandleRegistClusterMessage(t *testing.T) {
	//croot := newFakeRootClusterHandler(t)

	//err := croot.handleRegistClusterMessage()
}

type fakeCloudTunnel struct {
	broadcastCalled bool
	sendCalled      bool
}

func newFakeCloudTunnel() *fakeCloudTunnel {
	return &fakeCloudTunnel{
		broadcastCalled: false,
		sendCalled:      false,
	}
}

func (f *fakeCloudTunnel) reset() {
	f.broadcastCalled = false
	f.sendCalled = false
}

func (f *fakeCloudTunnel) Start() error {
	return nil
}

func (f *fakeCloudTunnel) Stop() error {
	return nil
}

func (f *fakeCloudTunnel) Send(clusterName string, msg []byte) error {
	f.sendCalled = true
	return nil
}

func (f *fakeCloudTunnel) Broadcast(msg []byte) {
	f.broadcastCalled = true
}

func (f *fakeCloudTunnel) RegistCheckNameValidFunc(fn tunnel.ClusterNameChecker) {}

func (f *fakeCloudTunnel) RegistAfterConnectHook(fn tunnel.AfterConnectHook) {}

func (f *fakeCloudTunnel) RegistReturnMessageFunc(fn tunnel.TunnelReadMessageFunc) {}

func (f *fakeCloudTunnel) RegistClientCloseHandler(fn tunnel.ClientCloseHandleFunc) {}

func newFakeRootClusterHandler(t *testing.T) *clusterHandler {
	ret := &clusterHandler{
		conf: &config.ClusterControllerConfig{
			ClusterUserDefineName: config.RootClusterName,
			TunnelListenAddr:      "8272",
			K8sClient:             oteclient.NewSimpleClientset(),
		},
		tunn:      fakeTunn,
		k8sEnable: false,
	}
	err := ret.valid()
	assert.Nil(t, err)
	fakeTunn.reset()
	return ret
}

func newFakeNoRootClusterHandler(t *testing.T) *clusterHandler {
	ret := &clusterHandler{
		conf: &config.ClusterControllerConfig{
			ClusterUserDefineName: "c1",
			TunnelListenAddr:      "8273",
			K8sClient:             oteclient.NewSimpleClientset(),
			ParentCluster:         "8272",
		},
		tunn:      fakeTunn,
		k8sEnable: false,
	}
	err := ret.valid()
	assert.Nil(t, err)
	fakeTunn.reset()
	return ret
}
