import { notification as antdNotification } from 'antd'
export const notification = {
  info(message: string, description: string) {
    antdNotification.info({
      message,
      description,
    })
  },
  error(message: string, description: string) {
    antdNotification.error({
      message,
      description,
    })
  },
}
