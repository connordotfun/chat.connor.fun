import { observable, action, reaction } from 'mobx'

class CommonStore {
  @observable appName = 'chat.connor.fun';
  @observable token = window.localStorage.getItem('jwt');
  @observable user = JSON.parse(window.localStorage.getItem('user'))
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

    reaction(
      () => this.user,
      user => {
        if (user) {
          window.localStorage.setItem('user', JSON.stringify(user))
        } else {
          window.localStorage.removeItem('user')
        }
      }
    )
  }

  @action setToken(token) {
    this.token = token
  }

  @action setUser(user) {
    this.user = user
  }

  @action setAppLoaded() {
    this.appLoaded = true
  }

}

export default new CommonStore()
