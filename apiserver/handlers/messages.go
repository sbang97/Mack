package handlers

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/info344-s17/challenges-sbang97/apiserver/models/messages"
	"github.com/info344-s17/challenges-sbang97/apiserver/models/users"
	"github.com/info344-s17/challenges-sbang97/apiserver/sessions"
)

func (ctx *Context) ChannelsHandler(w http.ResponseWriter, r *http.Request) {
	state := &SessionState{}
	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, state)
	if err != nil {
		http.Error(w, "failed to retrieve state:"+err.Error(), http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "GET":
		channels, err := ctx.MessageStore.GetAllChannels(state.User.ID)
		if err != nil {
			http.Error(w, "failed to retrieve user's channels", http.StatusBadRequest)
			return
		}
		w.Header().Add(headerContentType, contentTypeJSONUTF8)
		encoder := json.NewEncoder(w)
		encoder.Encode(channels)
	case "POST":
		newChannel := &messages.NewChannel{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(newChannel); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		chann, err := ctx.MessageStore.InsertChannel(newChannel, state.User.ID)
		if err != nil {
			http.Error(w, "error inserting the channel: "+err.Error(), http.StatusInternalServerError)
			return
		}
		channel := &ChannelEvent{
			Type:    newChannelCreated,
			Channel: chann,
		}
		ctx.Notifier.Notify(channel)

	default:
		http.Error(w, "request method must be GET or POST", http.StatusBadRequest)
		return
	}
}

func (ctx *Context) SpecificChannelHandler(w http.ResponseWriter, r *http.Request) {
	state := &SessionState{}
	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, state)
	if err != nil {
		http.Error(w, "failed to retrieve state:"+err.Error(), http.StatusBadRequest)
		return
	}
	url := r.URL
	_, path := path.Split(url.Path)
	chanID := messages.ChannelID(path)
	channel, err := ctx.MessageStore.GetChannelByID(chanID)
	switch r.Method {
	case "GET":
		messages, err := ctx.MessageStore.GetMessages(chanID, 500)
		if err != nil {
			http.Error(w, "error getting the messages:"+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add(headerContentType, contentTypeJSONUTF8)
		encoder := json.NewEncoder(w)
		encoder.Encode(messages)
	case "PATCH":
		updates := &messages.ChannelUpdates{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(updates); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if channel.CreatorID != state.User.ID {
			http.Error(w, "user is not the creator of this channel, cannot update:"+err.Error(), http.StatusBadRequest)
			return
		}
		updatedChann, err := ctx.MessageStore.UpdateChannel(updates, channel)
		if err != nil {
			http.Error(w, "error updating the channel:"+err.Error(), http.StatusBadRequest)
			return
		}
		updatedChannel := &ChannelEvent{
			Type:    channelUpdated,
			Channel: updatedChann,
		}
		ctx.Notifier.Notify(updatedChannel)
	case "DELETE":
		if channel.CreatorID != state.User.ID {
			http.Error(w, "user is not the creator of this channel, cannot delete:"+err.Error(), http.StatusBadRequest)
			return
		}
		err = ctx.MessageStore.DeleteChannel(chanID)
		if err != nil {
			http.Error(w, "error deleting the channel:"+err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	case "LINK":
		if channel.Private == true && state.User.ID != channel.CreatorID {
			http.Error(w, "The channel is private, user cannot be added to it:"+err.Error(), http.StatusBadRequest)
			return
		}
		userID := users.UserID(r.Header.Get("Link"))
		err = ctx.MessageStore.AddUser(userID, channel)
		if err != nil {
			http.Error(w, "failed to add user to members list:"+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case "UNLINK":
		if channel.Private == true && state.User.ID != channel.CreatorID {
			http.Error(w, "The channel is private, user cannot be removed:"+err.Error(), http.StatusBadRequest)
			return
		}
		userID := users.UserID(r.Header.Get("Link"))
		err = ctx.MessageStore.RemoveUser(userID, channel)
		if err != nil {
			http.Error(w, "failed to delete user from members list:"+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (ctx *Context) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	state := &SessionState{}
	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, state)
	if err != nil {
		http.Error(w, "failed to retrieve state:"+err.Error(), http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "POST":
		newMsg := &messages.NewMessage{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(newMsg); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		msg, err := ctx.MessageStore.InsertMessage(newMsg, state.User)
		if err != nil {
			http.Error(w, "failed to insert the new message:"+err.Error(), http.StatusBadRequest)
			return
		}
		msgEvent := &MessageEvent{
			Type:    newMessage,
			Message: msg,
		}
		ctx.Notifier.Notify(msgEvent)
	}
}

func (ctx *Context) SpecificMessageHandler(w http.ResponseWriter, r *http.Request) {
	state := &SessionState{}
	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, state)
	if err != nil {
		http.Error(w, "failed to retrieve state:"+err.Error(), http.StatusBadRequest)
		return
	}
	url := r.URL
	_, path := path.Split(url.Path)
	msgID := messages.MessageID(path)
	msg, _ := ctx.MessageStore.GetMessageByID(msgID)
	if state.User.ID != msg.CreatorID {
		http.Error(w, "user is not the creator of the message", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "PATCH":
		updates := &messages.MessageUpdates{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(updates); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		updatedMsg, err := ctx.MessageStore.UpdateMessage(updates, msg)
		if err != nil {
			http.Error(w, "error updating message:"+err.Error(), http.StatusInternalServerError)
			return
		}
		msgUpdate := &MessageEvent{
			Type:    messageUpdated,
			Message: updatedMsg,
		}
		ctx.Notifier.Notify(msgUpdate)
	case "DELETE":
		err = ctx.MessageStore.DeleteMessage(msgID)
		if err != nil {
			http.Error(w, "error deleting message:"+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
