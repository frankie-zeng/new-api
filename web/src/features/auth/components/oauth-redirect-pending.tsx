/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { Loader2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'

type OAuthRedirectPendingProps = {
  providerName?: string | null
}

export function OAuthRedirectPending(props: OAuthRedirectPendingProps) {
  const { t } = useTranslation()
  const headline = props.providerName
    ? t('Signing you in with {{provider}}', { provider: props.providerName })
    : t('Loading...')

  return (
    <div
      role='status'
      aria-live='polite'
      className='flex flex-col items-center space-y-4 text-center'
    >
      <Loader2 className='text-muted-foreground h-8 w-8 animate-spin' />
      <h2 className='text-2xl font-semibold tracking-tight'>{headline}</h2>
    </div>
  )
}
