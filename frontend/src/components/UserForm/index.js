import React, { Component } from 'react'
import './index.css'
import { inject, observer } from 'mobx-react'
import { withRouter } from 'react-router-dom'


@withRouter
@inject('authStore') @observer
class UserForm extends Component {
    componentWillUnmount() {
        this.props.authStore.reset()
    }
    handleUserChange = e => this.props.authStore.setUsername(e.target.value);
    handlePasswordChange = e => this.props.authStore.setPassword(e.target.value);
    handleSubmitForm = e => {e.preventDefault(); this.props.authStore.login() }
    render() {
        const { values, inProgress, errors } = this.props.authStore
        return (
            <div className="UserForm">
                <form onSubmit={this.handleSubmitForm}>
                    <h3>Log In</h3>
                    <fieldset className="form-group">
                        <input
                        className="form-control form-control-lg"
                        type="text"
                        placeholder="Username"
                        value={values.username}
                        onChange={this.handleUserChange}
                        />
                    </fieldset>

                    <fieldset className="form-group">
                        <input
                        className="form-control form-control-lg"
                        type="password"
                        placeholder="Password"
                        value={values.password}
                        onChange={this.handlePasswordChange}
                        />
                    </fieldset>

                    <button
                        className="convex"
                        type="submit"
                        disabled={inProgress}
                    >
                        Sign in
                    </button>
                </form>
                {errors ? <p>{errors}</p> : null}
            </div>
        )
    }
}

export default UserForm