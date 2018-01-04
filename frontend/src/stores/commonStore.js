import { observable, action, reaction } from 'mobx'

class CommonStore {
  @observable appName = 'chat.connor.fun';
  @observable token = window.localStorage.getItem('jwt');
  @observable username = undefined
  @observable appLoaded = false;

  constructor() {
    reaction(
      () => this.token,
      token => {
        if (token) {
          window.localStorage.setItem('jwt', token)
        } else {
          window.localStorage.removeItem('jwt')
        }
      }
    )
  }

  @action setToken(token) {
    this.token = token
  }

  @action setUsername(username) {
    this.username = username
  }

  @action setAppLoaded() {
    this.appLoaded = true
  }

}

export default new CommonStore()
