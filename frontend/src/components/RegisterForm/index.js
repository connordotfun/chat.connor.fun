import React, { Component } from 'react'
import './index.css'
import { inject, observer } from 'mobx-react'
import { Link } from 'react-router-dom'

@inject('authStore') @observer
class RegisterForm extends Component {
    componentWillUnmount() {
        this.props.authStore.reset()
    }
    handleUserChange = e => this.props.authStore.setUsername(e.target.value);
    handlePasswordChange = e => this.props.authStore.setPassword(e.target.value);
    handleEmailChange = e => this.props.authStore.setEmail(e.target.value)
    handleSubmitForm = e => {
        e.preventDefault()
        try {
            this.props.authStore.register()
        } catch (error) {
            console.log(error)
        }
    }
    render() {
        return (
            <div className="RegisterForm">
                <form onSubmit={this.handleSubmitForm}>
                    <h3>Make an account to start chatting.</h3>
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

                    <fieldset className="form-group">
                        <input
                        className="form-control form-control-lg"
                        type="email"
                        placeholder="Email"
                        value={this.props.authStore.values.email}
                        onChange={this.handleEmailChange}
                        />
                    </fieldset>

                    <button
                        type="submit"
                        disabled={this.props.authStore.inProgress}
                    >
                        Get Started
                    </button>
                    <p className="login-desc">Have an account? <Link to="/login">Log in</Link>.</p>
                    <Link to="/login"><button className="login-button">Log in</button></Link>
                </form>
                {this.props.authStore.errors ? <p>{this.props.authStore.errors}</p> : null}
            </div>
        )
    }
}

export default RegisterForm