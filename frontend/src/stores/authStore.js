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
  }

  @action register() {
    this.inProgress = true
    this.errors = undefined
  }

  @action logout() {
    commonStore.setToken(undefined)
    return new Promise(res => res())
  }
}

export default new AuthStore()
