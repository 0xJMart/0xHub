import { Route, Routes } from 'react-router-dom';
import { HubLayout } from '@/app/layouts/HubLayout';
import { ProjectDetailPage } from '@/features/projects/pages/ProjectDetailPage';
import { ProjectsListPage } from '@/features/projects/pages/ProjectsListPage';
import { NotFoundPage } from '@/features/system/NotFoundPage';

function App(): JSX.Element {
  return (
    <Routes>
      <Route element={<HubLayout />}>
        <Route index element={<ProjectsListPage />} />
        <Route path="projects/:slug" element={<ProjectDetailPage />} />
      </Route>
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}

export default App;
