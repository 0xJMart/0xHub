import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import App from '../App'
import { Project } from '../types'

// Mock the API module
vi.mock('../api', () => ({
  fetchProjects: vi.fn(),
}))

const { fetchProjects } = await import('../api')

describe('App', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should display loading state initially', () => {
    ;(fetchProjects as any).mockImplementation(() => new Promise(() => {})) // Never resolves

    render(<App />)

    // Check for loading spinner by class name
    const spinner = document.querySelector('.animate-spin')
    expect(spinner).toBeInTheDocument()
  })

  it('should display projects after loading', async () => {
    const mockProjects: Project[] = [
      {
        id: '1',
        name: 'Project 1',
        description: 'Description 1',
        url: 'https://project1.com',
        category: 'web',
      },
      {
        id: '2',
        name: 'Project 2',
        description: 'Description 2',
        url: 'https://project2.com',
        category: 'mobile',
      },
    ]

    ;(fetchProjects as any).mockResolvedValueOnce(mockProjects)

    render(<App />)

    await waitFor(() => {
      expect(screen.getByText('Project 1')).toBeInTheDocument()
      expect(screen.getByText('Project 2')).toBeInTheDocument()
    })

    // The text is split across elements, so we check for parts of it
    expect(screen.getByText(/Showing/i)).toBeInTheDocument()
    // Check that projects count is displayed (there are multiple "2"s, so we use getAllByText)
    const countElements = screen.getAllByText('2')
    expect(countElements.length).toBeGreaterThan(0)
  })

  it('should display error message when fetch fails', async () => {
    ;(fetchProjects as any).mockRejectedValueOnce(new Error('Failed to fetch'))

    render(<App />)

    await waitFor(() => {
      expect(screen.getByText(/Error loading projects/i)).toBeInTheDocument()
      expect(screen.getByText('Failed to fetch')).toBeInTheDocument()
    })
  })

  it('should display empty state when no projects', async () => {
    ;(fetchProjects as any).mockResolvedValueOnce([])

    render(<App />)

    await waitFor(() => {
      expect(screen.getByText(/No projects found/i)).toBeInTheDocument()
    })
  })
})

