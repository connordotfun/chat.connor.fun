import React, { Component } from 'react'
import './index.css'
import { inject, observer } from 'mobx-react'

@inject('authStore') @observer
class UserForm extends Component {
    componentWillUnmount() {
        this.props.authStore.reset()
    }
    handleUserChange = e => this.props.authStore.setUsername(e.target.value);
    handlePasswordChange = e => this.props.authStore.setPassword(e.target.value);
    handleSubmitForm = e => {e.preventDefault(); this.props.authStore.login() }
    render() {
        return (
            <div className="UserForm">
                <form onSubmit={this.handleSubmitForm}>
                    <h3>Log in to start chatting</h3>
                    <fieldset className="form-group">
                        <input
                        className="form-control form-control-lg"
                        type="text"
                        placeholder="Username"
                        value={this.props.authStore.values.username}
                        onChange={this.handleUserChange}
                        />
                    </fieldset>

                    <fieldset className="form-group">
                        <input
                        className="form-control form-control-lg"
                        type="password"
                        placeholder="Password"
                        value={this.props.authStore.values.password}
                        onChange={this.handlePasswordChange}
                        />
                    </fieldset>

                    <button
                        className="convex"
                        type="submit"
                        disabled={this.props.authStore.inProgress}
                    >
                        Log in
                    </button>
                </form>
                {this.props.authStore.errors ? <p>{this.props.authStore.errors}</p> : null}
            </div>
        )
    }
}

export default UserForm