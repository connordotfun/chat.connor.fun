import React, { Component } from 'react'
import './index.css'
import { inject, observer } from 'mobx-react'
import { withRouter } from 'react-router-dom'


@withRouter
@inject('authStore') @observer
class UserForm extends Component {
    constructor(props) {
        super(props)
    }
    componentWillUnmount() {
        this.props.authStore.reset()
    }
    handleUserChange = e => this.props.authStore.setUsername(e.target.value);
    handlePasswordChange = e => this.props.authStore.setPassword(e.target.value);
    handleSubmitForm = e => {e.preventDefault(); this.props.authStore.login() }
    render() {
        const { values, inProgress } = this.props.authStore
        return (
            <div className="UserForm">
                <form onSubmit={this.handleSubmitForm}>
                    <fieldset>

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
                        className="btn btn-lg btn-primary pull-xs-right"
                        type="submit"
                        disabled={inProgress}
                    >
                        Sign in
                    </button>

                    </fieldset>
                </form>
            </div>
        )
    }
}

export default UserForm