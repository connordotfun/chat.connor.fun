import { observable, action } from 'mobx'
import axios from 'axios'

class RoomStore {
    @observable roomId = undefined
    @observable messages = []
    @observable members = []

    @action
    joinRoom(roomName) {
        axios.get('/api/v1/rooms/' + roomName)
        .then(action((res) => {
            this.roomId = res.data.data.id
            this.members = res.data.data.members
            this.fetchMessages(50)
        }))
        .catch(action((e) => {
            // In the future, we may want to do something different here
            this.roomId = undefined
            this.messages = []
            this.members = []
        }))
    }

    @action
    leaveRoom() {
        this.roomId = undefined
        this.messages = []
        this.members = []
    }

    @action
    onMessage(evt) { // takes websocket onmessage data
        let message = JSON.parse(evt.data)[0]
        this.messages.push(message)
    }

    @action
    fetchMessages(count) {
        let params = {
            room_id: this.roomId
        }

        if (count) {
            params.count = count
        }

        if (this.roomId) {
            axios.get('/api/v1/messages', { params })
            .then(action((res) => {
                this.messages = res.data.data.sort(this.compareTime)
            }))
            .catch(action((e) => {
                this.messages = []
            }))
        } else {
            this.messages = []
        }
    }

    compareTime(a, b) {
        return a.createTime - b.createTime
    }
}

export default new RoomStore()
