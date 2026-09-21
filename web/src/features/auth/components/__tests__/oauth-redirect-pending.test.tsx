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
import { render, screen } from '@testing-library/react'
import { describe, expect, test } from 'vitest'

const { createInstance } = await import('i18next')
const { I18nextProvider, initReactI18next } = await import('react-i18next')
const { OAuthRedirectPending } = await import('../oauth-redirect-pending')

const i18n = createInstance()
await i18n.use(initReactI18next).init({
  lng: 'en',
  resources: {
    en: {
      translation: {
        'Loading...': 'Loading...',
        'Signing you in with {{provider}}': 'Signing you in with {{provider}}',
      },
    },
  },
})

function PendingHarness(props: { providerName?: string | null }) {
  return (
    <I18nextProvider i18n={i18n}>
      <OAuthRedirectPending providerName={props.providerName} />
    </I18nextProvider>
  )
}

describe('OAuth redirect pending screen', () => {
  test('announces a loading status when the provider is not known yet', () => {
    render(<PendingHarness />)

    const status = screen.getByRole('status')
    expect(status).toHaveAttribute('aria-live', 'polite')
    expect(status).toHaveTextContent('Loading...')
  })

  test('announces the provider name while redirecting to a known OAuth provider', () => {
    render(<PendingHarness providerName='Company SSO' />)

    expect(screen.getByRole('status')).toHaveTextContent(
      'Signing you in with Company SSO'
    )
  })
})
