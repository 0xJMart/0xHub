import { Project, ProjectsResponse } from './types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

export async function fetchProjects(): Promise<Project[]> {
  const response = await fetch(`${API_BASE_URL}/projects`);
  if (!response.ok) {
    throw new Error('Failed to fetch projects');
  }
  const data: ProjectsResponse = await response.json();
  return data.projects;
}

export async function fetchProject(id: string): Promise<Project> {
  const response = await fetch(`${API_BASE_URL}/projects/${id}`);
  if (!response.ok) {
    throw new Error('Failed to fetch project');
  }
  return response.json();
}

