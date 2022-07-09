import { notification } from 'services/notification'
import { ErrorResponse } from 'dto/common'

export function handleApiError(t: (k: string) => string, error: unknown) {
  const errorResponse = error as ErrorResponse
  if (errorResponse.message) {
    notification.error(t('action_error'), t(errorResponse.message))
  } else {
    notification.error(t('action_error'), t('unknown_error'))
  }
}
