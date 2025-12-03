import { describe, it, expect, vi, beforeEach } from 'vitest'
import { fetchProjects, fetchProject } from '../api'
import { Project } from '../types'

// Mock global fetch
global.fetch = vi.fn()

describe('API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('fetchProjects', () => {
    it('should fetch and return projects', async () => {
      const mockProjects: Project[] = [
        {
          id: '1',
          name: 'Test Project',
          description: 'A test project',
          url: 'https://test.com',
          category: 'testing',
        },
      ]

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ projects: mockProjects }),
      })

      const result = await fetchProjects()

      expect(global.fetch).toHaveBeenCalledWith('http://localhost:8080/api/projects')
      expect(result).toEqual(mockProjects)
    })

    it('should throw error when fetch fails', async () => {
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: false,
        status: 500,
      })

      await expect(fetchProjects()).rejects.toThrow('Failed to fetch projects')
    })

    it('should handle network errors', async () => {
      ;(global.fetch as any).mockRejectedValueOnce(new Error('Network error'))

      await expect(fetchProjects()).rejects.toThrow('Network error')
    })
  })

  describe('fetchProject', () => {
    it('should fetch and return a single project', async () => {
      const mockProject: Project = {
        id: '1',
        name: 'Test Project',
        description: 'A test project',
        url: 'https://test.com',
      }

      ;(global.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => mockProject,
      })

      const result = await fetchProject('1')

      expect(global.fetch).toHaveBeenCalledWith('http://localhost:8080/api/projects/1')
      expect(result).toEqual(mockProject)
    })

    it('should throw error when project not found', async () => {
      ;(global.fetch as any).mockResolvedValueOnce({
        ok: false,
        status: 404,
      })

      await expect(fetchProject('non-existent')).rejects.toThrow('Failed to fetch project')
    })
  })
})

