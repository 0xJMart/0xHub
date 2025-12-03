export interface Project {
  id: string;
  name: string;
  description: string;
  url: string;
  icon?: string;
  category?: string;
  status?: string;
}

export interface ProjectsResponse {
  projects: Project[];
}

