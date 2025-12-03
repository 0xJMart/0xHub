import { Project } from '../types'

interface ProjectCardProps {
  project: Project
}

export default function ProjectCard({ project }: ProjectCardProps) {
  return (
    <a
      href={project.url}
      target="_blank"
      rel="noopener noreferrer"
      className="block bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow duration-200 overflow-hidden border border-gray-200"
    >
      <div className="p-6">
        <div className="flex items-start space-x-4">
          {project.icon && (
            <img
              src={project.icon}
              alt={`${project.name} icon`}
              className="w-12 h-12 rounded-lg object-contain flex-shrink-0"
              onError={(e) => {
                // Hide image on error
                e.currentTarget.style.display = 'none'
              }}
            />
          )}
          <div className="flex-1 min-w-0">
            <h3 className="text-xl font-semibold text-gray-900 mb-2 truncate">
              {project.name}
            </h3>
            <p className="text-gray-600 text-sm mb-3 line-clamp-3">
              {project.description}
            </p>
            <div className="flex items-center space-x-3 flex-wrap gap-2">
              {project.category && (
                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                  {project.category}
                </span>
              )}
              {project.status && (
                <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                  project.status === 'active' 
                    ? 'bg-green-100 text-green-800' 
                    : 'bg-gray-100 text-gray-800'
                }`}>
                  {project.status}
                </span>
              )}
            </div>
          </div>
        </div>
      </div>
    </a>
  )
}

