import { observable, action } from 'mobx'
import axios from 'axios'
import commonStore from './commonStore'

class AuthStore {
  @observable inProgress = false;
  @observable errors = undefined;

  @observable values = {
    username: '',
    password: '',
  };

  @action setUsername(username) {
    this.values.username = username
  }

  @action setPassword(password) {
    this.values.password = password
  }

  @action reset() {
    this.values.username = ''
    this.values.password = ''
  }

  @action login() {
    this.inProgress = true
    this.errors = undefined
    axios.post('/api/v1/login', {
      username: this.values.username,
      secret: this.values.password
    })
    .then((res) => {
      commonStore.setToken(res.data.data.token)
    })
    .catch(action((err) => {
      this.errors = err.response.data.error
    }))
    .finally(action(() => {this.inProgress = false}))
  }

  @action register() {
    this.inProgress = true
    this.errors = undefined
    axios.post('/api/v1/users', {
      username: this.values.username,
      secret: this.values.password
    })
    .then((res) => {
      this.login()
    })
    .catch(action((err) => {
      this.errors = err.response.data.error
    }))
    .finally(action(() => {this.inProgress = false}))

  }

  @action logout() {
    commonStore.setToken(undefined)
    return new Promise(res => res())
  }
}

export default new AuthStore()
