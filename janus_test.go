package janus

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func Test_Connect(t *testing.T) {

	client, err := Connect("ws://localhost:8188/janus")
	if err != nil {
		t.Fail()
		return
	}
	// mess, err := client.Info()
	// if err != nil {
	// 	t.Fail()
	// 	return
	// }
	// t.Log(mess)

	id := uuid.NewString()
	mess := BaseMsg{
		Type: "create",
		ID:   id,
	}
	exp := &SuccessMsg{
		Type: "success",
		ID:   id,
	}

	sessMsg, err := client.CreateSession(context.Background(), mess)
	if err != nil {
		t.Fail()
		return
	}
	Assert(t, exp.Type, sessMsg.Type)
	Assert(t, exp.ID, sessMsg.ID)

	id = uuid.NewString()
	attachMsg := BaseMsg{
		Type:    "attach",
		ID:      id,
		Plugin:  "janus.plugin.echotest",
		Session: sessMsg.Data.ID,
	}
	sess := client.Sessions[sessMsg.Data.ID]

	pluginAttMsg, err := sess.AttachSession(context.Background(), attachMsg)
	if err != nil {
		t.Error(err)
	}

	exp = &SuccessMsg{
		Type: "success",
		ID:   id,
	}

	Assert(t, exp.Type, pluginAttMsg.Type)
	Assert(t, exp.ID, pluginAttMsg.ID)

	handler := sess.Handles[pluginAttMsg.Data.ID]
	id = uuid.NewString()
	msg := HandlerMessageJsep{
		HandlerMessage: HandlerMessage{
			BaseMsg: BaseMsg{
				Type:    "message",
				ID:      id,
				Session: sessMsg.Data.ID,
			},
			Body: map[string]any{
				"audio": true,
				"video": true,
			},
			Handle: pluginAttMsg.Data.ID,
		},
	}

	eventMsg, err := handler.Message(context.Background(), msg)
	if err != nil {
		t.Error(err)
	}
	expEvent := EventMsg{
		Type:   "ack",
		ID:     id,
		Handle: sessMsg.Data.ID,
		Plugindata: PluginData{
			Plugin: "janus.plugin.echotest",
			Data: map[string]interface{}{
				"echotest": "event",
				"result":   "ok",
			},
		},
	}
	if !reflect.DeepEqual(expEvent, *eventMsg) {
		t.Errorf("Expected <%+v>,got <+%v> instead", expEvent, *eventMsg)
	}
}

func Assert[T comparable](t *testing.T, exp T, got T) {
	t.Helper()
	if exp != got {
		t.Errorf("Expected %v, got %v instead", exp, got)
	}
}
