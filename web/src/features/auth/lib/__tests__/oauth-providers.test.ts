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
import { describe, expect, test } from 'vitest'

import type { CustomOAuthProviderInfo, SystemStatus } from '../../types'
import {
  getRedirectOAuthProviders,
  getSoleRedirectOAuthProvider,
  isForcedLocalAuth,
  shouldSkipLocalAuthForm,
} from '../oauth'

function customProvider(
  overrides: Partial<CustomOAuthProviderInfo> = {}
): CustomOAuthProviderInfo {
  return {
    id: 1,
    name: 'Company SSO',
    slug: 'company',
    icon: '',
    client_id: 'client-id',
    authorization_endpoint: 'https://sso.example.com/authorize',
    scopes: 'openid',
    ...overrides,
  }
}

describe('getRedirectOAuthProviders', () => {
  test('returns an empty list when status is missing', () => {
    expect(getRedirectOAuthProviders(null)).toEqual([])
  })

  test('ignores telegram and wechat because they cannot redirect', () => {
    const status: SystemStatus = {
      telegram_oauth: true,
      telegram_bot_name: 'bot',
      wechat_login: true,
    }

    expect(getRedirectOAuthProviders(status)).toEqual([])
  })

  test('ignores a provider that is enabled but missing its client id', () => {
    const status: SystemStatus = {
      github_oauth: true,
      discord_oauth: true,
      discord_client_id: '',
    }

    expect(getRedirectOAuthProviders(status)).toEqual([])
  })

  test('returns configured redirect providers in a stable order', () => {
    const status: SystemStatus = {
      github_oauth: true,
      github_client_id: 'gh',
      oidc_enabled: true,
      oidc_client_id: 'oidc',
      oidc_authorization_endpoint: 'https://idp.example.com/authorize',
      oidc_display_name: 'Acme ID',
      custom_oauth_providers: [customProvider()],
    }

    expect(
      getRedirectOAuthProviders(status).map((provider) => provider.kind)
    ).toEqual(['github', 'oidc', 'custom'])
  })
})

describe('getSoleRedirectOAuthProvider', () => {
  test('returns the only redirect provider when exactly one is configured', () => {
    const status: SystemStatus = {
      custom_oauth_providers: [customProvider()],
    }

    expect(getSoleRedirectOAuthProvider(status)).toEqual({
      kind: 'custom',
      name: 'Company SSO',
      provider: customProvider(),
    })
  })

  test('returns null when more than one redirect provider is configured', () => {
    const status: SystemStatus = {
      github_oauth: true,
      github_client_id: 'gh',
      linuxdo_oauth: true,
      linuxdo_client_id: 'ld',
    }

    expect(getSoleRedirectOAuthProvider(status)).toBeNull()
  })
})

describe('shouldSkipLocalAuthForm', () => {
  test('skips the local form once a redirect OAuth provider is configured', () => {
    const status: SystemStatus = {
      github_oauth: true,
      github_client_id: 'gh',
    }

    expect(shouldSkipLocalAuthForm(status)).toBe(true)
  })

  test('keeps the local form when forceLocalAuth is set', () => {
    const status: SystemStatus = {
      github_oauth: true,
      github_client_id: 'gh',
    }

    expect(shouldSkipLocalAuthForm(status, { forceLocalAuth: true })).toBe(
      false
    )
  })

  test('keeps the local form when no redirect OAuth provider is ready', () => {
    expect(shouldSkipLocalAuthForm({ telegram_oauth: true })).toBe(false)
  })
})

describe('isForcedLocalAuth', () => {
  test('accepts boolean and string query values', () => {
    expect(isForcedLocalAuth(true)).toBe(true)
    expect(isForcedLocalAuth('true')).toBe(true)
    expect(isForcedLocalAuth('1')).toBe(true)
    expect(isForcedLocalAuth(false)).toBe(false)
    expect(isForcedLocalAuth('false')).toBe(false)
    expect(isForcedLocalAuth(undefined)).toBe(false)
  })
})
