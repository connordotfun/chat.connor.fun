import { observable, action } from 'mobx'
import commonStore from './commonStore'

class SocketStore {
  @observable socket = undefined
  @observable connected = false
  @observable error = undefined
  @observable listeners = []

  @observable message = ""

  @action joinRoom(room) {
    if (this.socket) {
      this.leaveRoom()
    }

    this.socket = new WebSocket((window.location.protocol === "https:"? "wss://" : "ws://") + window.location.host + "/api/v1/rooms/" + room +  "/messages/ws", commonStore.token)
    this.socket.onopen = (e) => { this.connected = true }
    this.socket.onerror = this.setError
    this.socket.onmessage = (e) => {
      this.listeners.map((fxn) => fxn(e))
    }
  }

  @action addListener(fxn) {
    this.listeners.push(fxn)
  }

  @action resetListeners() {
    this.listeners = []
  }

  @action leaveRoom() {
    if (this.socket) {
      this.socket.close()
      this.socket = undefined
    }

    this.connected = false
  }

  @action setMessage(message) {
    this.message = message
  }

  @action sendMessage() {
    if (this.message !== "") {
      this.socket.send(JSON.stringify({
        text: this.message
      }))
    }
  }

  @action setError(e) {
    this.error = e
  }
}

export default new SocketStore()
