import { describe, it, expect, vi, beforeEach } from 'vitest'
import router from '@/router'
import axios from 'axios'

// Mock axios requests
vi.mock('axios')

describe('Router Guards (beforeEnter)', () => {

  beforeEach(() => {
    vi.clearAllMocks()
  })

  // PUBLIC ROUTES

  it('allows access to home without authentication', async () => {
    await router.push('/')
    await router.isReady()
    // check 1
    expect(router.currentRoute.value.name).toBe('home')
  })

  it('allows access to login without authentication', async () => {
    await router.push('/login')
    await router.isReady()

    // check 2
    expect(router.currentRoute.value.name).toBe('login')
  })

  it('allows access to signup without authentication', async () => {
    await router.push('/signup')
    await router.isReady()

    //check 3
    expect(router.currentRoute.value.name).toBe('signup')
  })

  // PROTECTED ROUTES (UNAUTHENTICATED)

  it('redirects to login when not authenticated', async () => {
    axios.get.mockResolvedValue({
      data: { authenticated: false }
    })

    await router.push('/dashboard')
    await router.isReady()

    expect(router.currentRoute.value.name).toBe('login')
  })

  // PROTECTED ROUTES (AUTHENTICATED)

  it('allows access to protected route when authenticated', async () => {
    axios.get.mockResolvedValue({
      data: { authenticated: true }
    })

    await router.push('/dashboard')
    await router.isReady()

    expect(router.currentRoute.value.name).toBe('dashboard')
  })

  // MULTIPLE ROUTES TEST

  const protectedRoutes = [
    '/library',
    '/profile',
    '/settings',
    '/dashboard',
    '/jobs/create'
  ]

  protectedRoutes.forEach((route) => {
    it(`blocks ${route} when unauthenticated`, async () => {
      axios.get.mockResolvedValue({
        data: { authenticated: false }
      })

      await router.push(route)
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('login')
    })
  })

  // API FAILURE CASE

  it('redirects to login if auth API fails', async () => {
    axios.get.mockRejectedValue(new Error('API down'))

    await router.push('/dashboard')
    await router.isReady()

    expect(router.currentRoute.value.name).toBe('login')
  })

})
