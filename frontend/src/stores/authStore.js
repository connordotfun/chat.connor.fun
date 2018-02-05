import { observable, action } from 'mobx'
import axios from 'axios'
import commonStore from './commonStore'

class AuthStore {
  @observable inProgress = false;
  @observable errors = undefined;

  @observable values = {
    username: '',
    password: '',
    email: ''
  };

  @action setUsername(username) {
    this.values.username = username
  }

  @action setPassword(password) {
    this.values.password = password
  }

  @action setEmail(email) {
    this.values.email = email
  }

  @action reset() {
    this.values.username = ''
    this.values.password = ''
    this.values.email = ''
  }

  @action login() {
    this.inProgress = true
    this.errors = undefined
    axios.post('/api/v1/login', {
      username: this.values.username,
      secret: this.values.password
    })
    .then((res) => {
      commonStore.setUser(res.data.data.user)
      commonStore.setToken(res.data.data.token)
    })
    .catch(action((err) => {
      this.errors = err.response.data.error
    }))
    .finally(action(() => {
      this.inProgress = false
    }))
  }

  @action register() {
    this.inProgress = true
    this.errors = undefined
    axios.post('/api/v1/users', {
      username: this.values.username,
      secret: this.values.password,
      email: this.values.email
    })
    .catch(action((err) => {
      this.errors = err.response.data.error
      throw this.errors
    }))
    .finally(action(() => {
      this.login()
      this.inProgress = false
    }))
  }

  @action verifyAccount(code) {
    this.inProgress = true
    this.errors = undefined
    axios.put('/api/v1/verifications/accountverification', {code}, {
      headers: {'Authorization': 'Bearer ' + commonStore.token}
    })
    .catch(action((err) => {
      this.errors = err.response.data.error
      throw this.errors
    }))
    .finally(action(() => {
      console.log('getting user')
      this.getUser()
      this.inProgress = false
    }))
  }

  @action getUser() {
    this.inProgress = true
    this.errors = undefined
    axios.get('/api/v1/users/' + commonStore.user.id, {
      headers: {'Authorization': 'Bearer ' + commonStore.token}
    })
    .catch(action((err) => {
      this.errors = err.response.data.error
      throw this.errors
    }))
    .then(action((res) => {
      commonStore.setUser(res.data.data)
    }))
    .finally(action(() => {
      this.inProgress = false
    }))
  }

  @action logout() {
    commonStore.setToken(undefined)
    commonStore.setUser(undefined)
    return new Promise(res => res())
  }
}

export default new AuthStore()
