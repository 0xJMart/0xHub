import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import ProjectCard from '../ProjectCard'
import { Project } from '../../types'

describe('ProjectCard', () => {
  const mockProject: Project = {
    id: '1',
    name: 'Test Project',
    description: 'This is a test project description',
    url: 'https://test.com',
    category: 'testing',
    status: 'active',
  }

  it('should render project name and description', () => {
    render(<ProjectCard project={mockProject} />)

    expect(screen.getByText('Test Project')).toBeInTheDocument()
    expect(screen.getByText('This is a test project description')).toBeInTheDocument()
  })

  it('should render project link with correct URL', () => {
    render(<ProjectCard project={mockProject} />)

    const link = screen.getByRole('link')
    expect(link).toHaveAttribute('href', 'https://test.com')
    expect(link).toHaveAttribute('target', '_blank')
    expect(link).toHaveAttribute('rel', 'noopener noreferrer')
  })

  it('should render category badge when provided', () => {
    render(<ProjectCard project={mockProject} />)

    expect(screen.getByText('testing')).toBeInTheDocument()
  })

  it('should render status badge when provided', () => {
    render(<ProjectCard project={mockProject} />)

    expect(screen.getByText('active')).toBeInTheDocument()
  })

  it('should not render category badge when not provided', () => {
    const projectWithoutCategory = { ...mockProject, category: undefined }
    render(<ProjectCard project={projectWithoutCategory} />)

    expect(screen.queryByText('testing')).not.toBeInTheDocument()
  })

  it('should not render status badge when not provided', () => {
    const projectWithoutStatus = { ...mockProject, status: undefined }
    render(<ProjectCard project={projectWithoutStatus} />)

    expect(screen.queryByText('active')).not.toBeInTheDocument()
  })

  it('should render icon when provided', () => {
    const projectWithIcon = { ...mockProject, icon: 'https://test.com/icon.png' }
    render(<ProjectCard project={projectWithIcon} />)

    const icon = screen.getByAltText('Test Project icon')
    expect(icon).toBeInTheDocument()
    expect(icon).toHaveAttribute('src', 'https://test.com/icon.png')
  })

  it('should not render icon when not provided', () => {
    const projectWithoutIcon = { ...mockProject, icon: undefined }
    const { container } = render(<ProjectCard project={projectWithoutIcon} />)

    const images = container.querySelectorAll('img')
    expect(images.length).toBe(0)
  })
})

