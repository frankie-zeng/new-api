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
import { useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { clearAuthentication } from '@/lib/api'
import {
  buildDiscordOAuthUrl,
  buildGitHubOAuthUrl,
  buildLinuxDOOAuthUrl,
  buildOIDCOAuthUrl,
} from '@/lib/oauth'

import { createOAuthFlow, logout } from '../api'
import {
  getSoleRedirectOAuthProvider,
  type RedirectOAuthProvider,
} from '../lib/oauth'
import { rememberOAuthLoginRedirect } from '../lib/oauth-callback-mode'
import type { SystemStatus } from '../types'

type UseAutoOAuthRedirectOptions = {
  enabled?: boolean
  redirectTo?: string
}

async function startRedirectOAuthLogin(
  status: SystemStatus,
  provider: RedirectOAuthProvider,
  redirectTo?: string
) {
  const response = await logout()
  if (!response.success) {
    throw new Error(response.message || 'Failed to sign out session')
  }
  clearAuthentication()

  const providerKey =
    provider.kind === 'custom' ? provider.provider.slug : provider.kind
  const state = await createOAuthFlow(providerKey, 'login')
  rememberOAuthLoginRedirect(state, redirectTo)
  const url = buildAutoOAuthUrl(provider, status, state)
  if (!url) {
    throw new Error('OAuth provider is not ready')
  }
  window.location.assign(url)
}

function buildAutoOAuthUrl(
  provider: RedirectOAuthProvider,
  status: SystemStatus,
  state: string
): string | null {
  switch (provider.kind) {
    case 'github':
      if (!status.github_client_id) return null
      return buildGitHubOAuthUrl(status.github_client_id, state)
    case 'discord':
      if (!status.discord_client_id) return null
      return buildDiscordOAuthUrl(status.discord_client_id, state)
    case 'oidc':
      if (!status.oidc_authorization_endpoint || !status.oidc_client_id) {
        return null
      }
      return buildOIDCOAuthUrl(
        status.oidc_authorization_endpoint,
        status.oidc_client_id,
        state
      )
    case 'linuxdo':
      if (!status.linuxdo_client_id) return null
      return buildLinuxDOOAuthUrl(status.linuxdo_client_id, state)
    case 'custom': {
      if (
        !provider.provider.authorization_endpoint ||
        !provider.provider.client_id
      ) {
        return null
      }
      const url = new URL(provider.provider.authorization_endpoint)
      url.searchParams.set('client_id', provider.provider.client_id)
      url.searchParams.set(
        'redirect_uri',
        `${window.location.origin}/oauth/${provider.provider.slug}`
      )
      url.searchParams.set('response_type', 'code')
      url.searchParams.set('state', state)
      if (provider.provider.scopes) {
        url.searchParams.set('scope', provider.provider.scopes)
      }
      return url.toString()
    }
  }
}

/**
 * When exactly one redirect-capable OAuth provider is configured, start it
 * immediately instead of rendering the local sign-in / sign-up form.
 */
export function useAutoOAuthRedirect(
  status: SystemStatus | null,
  options?: UseAutoOAuthRedirectOptions
) {
  const { t } = useTranslation()
  const enabled = options?.enabled !== false
  const redirectTo = options?.redirectTo

  const provider = useMemo(() => {
    if (!enabled) return null
    return getSoleRedirectOAuthProvider(status)
  }, [enabled, status])

  const [failed, setFailed] = useState(false)

  useEffect(() => {
    if (!provider || !status) {
      setFailed(false)
      return
    }

    let cancelled = false
    void startRedirectOAuthLogin(status, provider, redirectTo).catch(() => {
      if (cancelled) return
      toast.error(
        t('Failed to start {{provider}} login', { provider: provider.name })
      )
      setFailed(true)
    })

    return () => {
      cancelled = true
    }
  }, [provider, status, redirectTo, t])

  return {
    isRedirecting: provider !== null && !failed,
    providerName: provider?.name ?? null,
  }
}
