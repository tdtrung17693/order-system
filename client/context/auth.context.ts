import React from 'react'
import { User } from '../dto/user.dto'
import { Maybe } from '../types/maybe'

export const AuthContext = React.createContext<{
  authenticated: boolean
  user: Maybe<User>
  logout: () => void
}>({
  authenticated: false,
  user: null,
  logout: () => {},
})
