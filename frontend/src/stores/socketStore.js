import { observable, action } from 'mobx'
import commonStore from './commonStore'
import roomStore from './roomStore'

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

    roomStore.joinRoom(room)
    this.socket = new WebSocket((window.location.protocol === "https:"? "wss://" : "ws://") + window.location.host + "/api/v1/rooms/" + room +  "/ws", commonStore.token)
    this.socket.onopen = action((e) => { this.connected = true })
    this.socket.onerror = this.setError
    this.socket.onclose = action((e) => {
      this.connected = false
    })
    this.socket.onmessage = roomStore.onMessage.bind(roomStore)
    // maybe add this back in if needed
    // (e) => {
    //   this.listeners.map((fxn) => fxn(e))
    // }
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
    roomStore.leaveRoom()
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
