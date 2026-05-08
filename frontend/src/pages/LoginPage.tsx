import { FormEvent, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'
import { useToast } from '../components/ToastProvider'

export const LoginPage = () => {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [credentials, setCredentials] = useState({ username: '', password: '' })
  const { pushToast } = useToast()

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault()
    try {
      await login(credentials.username, credentials.password)
      pushToast('success', 'Signed in successfully.')
      navigate('/dashboard')
    } catch (err) {
      pushToast('error', (err as Error).message)
    }
  }

  return (
    <div className="login-page">
      <section className="login-card">
        <div className="login-hero">
          <span className="kicker">OCS Access</span>
          <h1>Secure room scheduling</h1>
          <p>
            Sign in with admin-issued credentials to book rooms and review
            upcoming schedules.
          </p>
        </div>
        <form className="form" onSubmit={handleSubmit}>
          <label className="field">
            <span>Username</span>
            <input
              className="input"
              value={credentials.username}
              onChange={(event) =>
                setCredentials({ ...credentials, username: event.target.value })
              }
              placeholder="ocs_admin"
              required
            />
          </label>
          <label className="field">
            <span>Password</span>
            <input
              className="input"
              type="password"
              value={credentials.password}
              onChange={(event) =>
                setCredentials({ ...credentials, password: event.target.value })
              }
              required
            />
          </label>
          <button className="button" type="submit">
            Sign in
          </button>
        </form>
      </section>
    </div>
  )
}
