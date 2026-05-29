import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { 
  listBatchesByCourse, 
  createBatchForCourse, 
  getBatchDetails, 
  listCoursesBySession, 
  assignCourseToSession,
  listCoursesForUser,
} from "../api/manageCompany.js";

const demo = {
  sessions: [
    { id: 'S-2025', name: '2025 - 2026', status: 'Active', start: '01 Jun, 2025', end: '31 May, 2026' },
    { id: 'S-2024', name: '2024 - 2025', status: 'Active', start: '01 Jun, 2024', end: '31 May, 2025' },
    { id: 'S-2023', name: '2023 - 2024', status: 'Completed', start: '01 Jun, 2023', end: '31 May, 2024' },
  ],
  courses: [
    { id: 'C-201', title: 'Full Stack Web Development', level: 'Advanced', seats: '28/40', dates: '01 Jun, 2025 - 30 Nov, 2025', status: 'Active' },
    { id: 'C-314', title: 'UI/UX Design Fundamentals', level: 'Intermediate', seats: '18/30', dates: '15 Jun, 2025 - 15 Dec, 2025', status: 'Active' },
    { id: 'C-487', title: 'Digital Marketing Mastery', level: 'Beginner', seats: '22/40', dates: '01 Jul, 2025 - 31 Dec, 2025', status: 'Active' },
  ],
  batches: [
    { id: 'B-2025F', name: 'Batch - Fall 2025', course: 'Full Stack Web Dev', start: '01 Jun, 2025', end: '30 Nov, 2025', students: '28/40', status: 'Active' },
    { id: 'B-2025W', name: 'Batch - Winter 2025', course: 'UI/UX Design', start: '15 Jun, 2025', end: '15 Dec, 2025', students: '18/30', status: 'Active' },
    { id: 'B-2025S', name: 'Batch - Summer 2025', course: 'Digital Marketing', start: '01 Jul, 2025', end: '31 Dec, 2025', students: '22/40', status: 'Active' },
  ],
  team: [
    { id: 'TM1', name: 'John Doe', role: 'Academic Head', dept: 'Academics', status: 'Active' },
    { id: 'TM2', name: 'Sarah Smith', role: 'Course Manager', dept: 'Academics', status: 'Active' },
    { id: 'TM3', name: 'Mike Johnson', role: 'Admissions Lead', dept: 'Admissions', status: 'Active' },
    { id: 'TM4', name: 'Emily Davis', role: 'Marketing Manager', dept: 'Marketing', status: 'Inactive' },
  ],
};

export default function ManageCompanyBatches() {
  const navigate = useNavigate();
  const [batches, setBatches] = useState([]);
  const [courses, setCourses] = useState([]);
  const [sessionCourses, setSessionCourses] = useState({});
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const token = localStorage.getItem('auth_token');

  const handleLogout = () => {
    localStorage.removeItem('auth_token');
    navigate('/');
  };

  useEffect(() => {
    const fetchCourses = async () => {
      if (!token) return;
      try {
        setLoading(true);
        const data = await listCoursesForUser({ page: 1, size: 20 }, token);
        setCourses(data.courses || []);
        // Auto-load batches for first course
        if (data.courses && data.courses.length > 0) {
          const courseData = await listBatchesByCourse(data.courses[0].id, { page: 1, size: 20 }, token);
          setBatches(courseData.batches || []);
        }
      } catch (err) {
        console.error('Failed to fetch courses for user', err);
      } finally {
        setLoading(false);
      }
    };
    fetchCourses();
  }, [token]);

  // Fetch batches by course
  const handleFetchBatchesByCourse = async (courseId) => {
    try {
      setLoading(true);
      const data = await listBatchesByCourse(courseId, { page: 1, size: 20 }, token);
      setBatches(data.batches || []);
      setError(null);
    } catch (err) {
      setError('Failed to fetch batches: ' + err.message);
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // Create new batch
  const handleCreateBatch = async (courseId) => {
    const batchName = prompt('Enter batch name:');
    if (!batchName) return;

    const startDate = prompt('Enter start date (YYYY-MM-DD):');
    if (!startDate) return;

    const endDate = prompt('Enter end date (YYYY-MM-DD):');
    if (!endDate) return;

    const maxStudents = prompt('Enter max students:');
    if (!maxStudents) return;

    try {
      setLoading(true);
      const payload = {
        batch_name: batchName,
        start_date: startDate,
        end_date: endDate,
        max_students: parseInt(maxStudents),
        status: 'draft',
      };
      const data = await createBatchForCourse(courseId, payload, token);
      alert('Batch created successfully');
      setBatches([...batches, data]);
      setError(null);
    } catch (err) {
      setError('Failed to create batch: ' + err.message);
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // Get batch details
  const handleGetBatchDetails = async (courseId, batchId) => {
    try {
      setLoading(true);
      const data = await getBatchDetails(courseId, batchId, token);
      console.log('Batch details:', data);
      alert('Batch details: ' + JSON.stringify(data, null, 2));
      setError(null);
    } catch (err) {
      setError('Failed to fetch batch details: ' + err.message);
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // List courses by session
  const handleListCoursesBySession = async (sessionId) => {
    try {
      setLoading(true);
      const data = await listCoursesBySession(sessionId, { page: 1, size: 20 }, token);
      setSessionCourses({ ...sessionCourses, [sessionId]: data.courses || [] });
      console.log('Courses in session:', data);
      setError(null);
    } catch (err) {
      setError('Failed to fetch courses: ' + err.message);
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // Assign course to session
  const handleAssignCourseToSession = async (sessionId, courseId) => {
    try {
      setLoading(true);
      const data = await assignCourseToSession(sessionId, courseId, token);
      alert('Course assigned to session: ' + data.course_id);
      await handleListCoursesBySession(sessionId);
      setError(null);
    } catch (err) {
      setError('Failed to assign course: ' + err.message);
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

    return (
      <div className="min-h-screen w-full bg-slate-50">
        <div className="flex">
          
          <aside className="hidden w-64 shrink-0 border-r border-slate-100 bg-white lg:block">
            
            <nav className="px-3 pb-6">
              <ul className="space-y-1 text-sm text-slate-700">
                <li className="px-3 py-2 rounded-lg hover:bg-slate-100">Dashboard</li>
                <li className="px-3 py-2 rounded-lg bg-indigo-50 text-indigo-700 font-semibold">Manage Company</li>
                <li className="px-3 py-2 rounded-lg hover:bg-slate-100">Academic Sessions</li>
                <li className="px-3 py-2 rounded-lg hover:bg-slate-100">Courses / Packages</li>
                <li className="px-3 py-2 rounded-lg hover:bg-slate-100">Batches</li>
                <li className="px-3 py-2 rounded-lg hover:bg-slate-100">Internal Team</li>
              </ul>
              <div className="mt-8 pt-6 border-t border-slate-200">
                <button
                  onClick={handleLogout}
                  className="w-full px-3 py-2 rounded-lg text-red-600 hover:bg-red-50 text-sm font-medium text-left"
                >
                  Logout
                </button>
              </div>
            </nav>
          </aside>

        <main className="flex-1 p-6">
          <header className="mb-6 flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-semibold text-slate-900">Manage Company</h1>
              <p className="text-sm text-slate-500">Manage academic sessions, courses, batches and internal team.</p>
            </div>
            <div className="flex items-center gap-4">
              <div className="hidden items-center gap-2 rounded border border-slate-100 bg-white px-3 py-2 sm:flex">
                <svg width="16" height="16" fill="none" className="text-slate-400"><circle cx="8" cy="8" r="7" stroke="#E5E7EB"/></svg>
                <input className="w-80 text-sm outline-none" placeholder="Search anything..." />
              </div>
              <button className="rounded bg-indigo-600 px-4 py-2 text-white">+ Academic Session</button>
              {error && <div className="text-red-600 text-sm">{error}</div>}
              {loading && <div className="text-blue-600 text-sm">Loading...</div>}
            </div>
          </header>

          <section className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4 mb-6">
            <div className="rounded-xl border border-slate-100 bg-white p-4">
              <p className="text-xs text-slate-500">Academic Sessions</p>
              <p className="mt-2 text-2xl font-semibold">3</p>
              <p className="text-sm text-slate-400">2 upcoming</p>
            </div>
            <div className="rounded-xl border border-slate-100 bg-white p-4">
              <p className="text-xs text-slate-500">Courses / Packages</p>
              <p className="mt-2 text-2xl font-semibold">24</p>
              <p className="text-sm text-slate-400">3 draft</p>
            </div>
            <div className="rounded-xl border border-slate-100 bg-white p-4">
              <p className="text-xs text-slate-500">Batches</p>
              <p className="mt-2 text-2xl font-semibold">18</p>
              <p className="text-sm text-slate-400">5 active</p>
            </div>
            <div className="rounded-xl border border-slate-100 bg-white p-4">
              <p className="text-xs text-slate-500">Team Members</p>
              <p className="mt-2 text-2xl font-semibold">12</p>
              <p className="text-sm text-slate-400">2 inactive</p>
            </div>
          </section>

          <section className="grid gap-6 lg:grid-cols-3">
            <div className="lg:col-span-2 space-y-6">
              <div className="rounded-xl border border-slate-100 bg-white p-4">
                <div className="flex items-center justify-between">
                  <h3 className="text-sm font-semibold text-slate-900">Academic Sessions</h3>
                  <a className="text-sm text-indigo-600 cursor-pointer">View all</a>
                </div>
                <div className="mt-4 overflow-x-auto">
                  <table className="min-w-full text-left text-sm">
                    <thead className="bg-slate-50 text-xs uppercase text-slate-500">
                      <tr>
                        <th className="px-4 py-3">Session Name</th>
                        <th className="px-4 py-3">Status</th>
                        <th className="px-4 py-3">Start Date</th>
                        <th className="px-4 py-3">End Date</th>
                        <th className="px-4 py-3">Actions</th>
                      </tr>
                    </thead>
                    <tbody>
                      {demo.sessions.map((s) => (
                        <tr key={s.id} className="border-t border-slate-100">
                          <td className="px-4 py-3">{s.name}</td>
                          <td className="px-4 py-3"><span className={`rounded-full px-2 py-1 text-xs ${s.status==='Active'? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-600'}`}>{s.status}</span></td>
                          <td className="px-4 py-3">{s.start}</td>
                          <td className="px-4 py-3">{s.end}</td>
                          <td className="px-4 py-3">
                            <button 
                              onClick={() => handleListCoursesBySession(s.id)}
                              className="text-indigo-600 hover:text-indigo-700 text-sm"
                            >
                              View Courses
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
                <p className="mt-3 text-xs text-slate-500">Total 3 sessions</p>
              </div>

              <div className="rounded-xl border border-slate-100 bg-white p-4">
                <div className="flex items-center justify-between">
                  <h3 className="text-sm font-semibold text-slate-900">Batches</h3>
                  <div className="flex items-center gap-3">
                    <a 
                      onClick={() => handleFetchBatchesByCourse(courses[0]?.id || demo.courses[0]?.id)}
                      className="text-sm text-indigo-600 cursor-pointer"
                    >
                      View all
                    </a>
                    <button 
                      onClick={() => handleCreateBatch(courses[0]?.id || demo.courses[0]?.id)}
                      className="rounded bg-white px-3 py-1 text-sm border border-slate-200"
                    >
                      + Add Batch
                    </button>
                  </div>
                </div>
                <div className="mt-4 overflow-x-auto">
                  <table className="min-w-full text-left text-sm">
                    <thead className="bg-slate-50 text-xs uppercase text-slate-500">
                      <tr>
                        <th className="px-4 py-3">Batch Name</th>
                        <th className="px-4 py-3">Course</th>
                        <th className="px-4 py-3">Start Date</th>
                        <th className="px-4 py-3">End Date</th>
                        <th className="px-4 py-3">Students</th>
                        <th className="px-4 py-3">Status</th>
                        <th className="px-4 py-3">Actions</th>
                      </tr>
                    </thead>
                    <tbody>
                      {batches.length > 0 ? (
                        batches.map((b, index) => (
                          <tr key={b.id || b.ID || `${b.course_id || b.CourseID || "course"}-${b.batch_name || b.BatchName || index}-${index}`} className="border-t border-slate-100">
                            <td className="px-4 py-3">{b.batch_name || b.name}</td>
                            <td className="px-4 py-3">{(courses.find && courses.find(c => c.id === b.course_id)) ? (courses.find(c => c.id === b.course_id).title || courses.find(c => c.id === b.course_id).Title) : b.course || b.course_id}</td>
                            <td className="px-4 py-3">{(b.start_date && new Date(b.start_date).toLocaleDateString()) || b.start}</td>
                            <td className="px-4 py-3">{(b.end_date && new Date(b.end_date).toLocaleDateString()) || b.end}</td>
                            <td className="px-4 py-3">{b.max_students || b.MaxStudents || b.students}</td>
                            <td className="px-4 py-3"><span className="rounded-full bg-emerald-50 px-2 py-1 text-xs text-emerald-700">{b.status || b.Status}</span></td>
                            <td className="px-4 py-3">
                              <button 
                                onClick={() => handleGetBatchDetails(courses[0]?.id, b.id)}
                                className="text-indigo-600 hover:text-indigo-700 text-sm"
                              >
                                View
                              </button>
                            </td>
                          </tr>
                        ))
                      ) : (
                        <tr>
                          <td colSpan="7" className="px-4 py-8 text-center text-slate-500">
                            <p className="text-sm">No batches available</p>
                          </td>
                        </tr>
                      )}
                    </tbody>
                  </table>
                </div>
                <p className="mt-3 text-xs text-slate-500">Total {batches.length} batches</p>
              </div>
            </div>

            <div className="space-y-6">
              <div className="rounded-xl border border-slate-100 bg-white p-4">
                <div className="flex items-center justify-between">
                  <h3 className="text-sm font-semibold text-slate-900">Courses / Packages</h3>
                  <div className="flex items-center gap-3">
                    <a className="text-sm text-indigo-600">View all</a>
                    <button className="rounded bg-white px-3 py-1 text-sm border border-slate-200">+ Add Course</button>
                  </div>
                </div>
                <div className="mt-4 grid gap-3">
                    {(courses.length ? courses : demo.courses).map((c) => (
                    <div key={c.id} className="rounded-xl border border-slate-100 bg-slate-50 p-3">
                      <div className="flex items-start justify-between">
                        <div>
                          <p className="text-xs text-slate-500">{c.level}</p>
                          <p className="mt-1 font-semibold text-slate-900">{c.title || c.title}</p>
                        </div>
                        <div className="text-right">
                          <p className="text-sm text-slate-700">{c.booked_seats ? `${c.booked_seats}/${c.total_seats}` : c.seats}</p>
                        </div>
                      </div>
                      <p className="mt-2 text-xs text-slate-500">{c.start_date ? new Date(c.start_date).toLocaleDateString() : c.dates}</p>
                    </div>
                  ))}
                </div>
              </div>

              <div className="rounded-xl border border-slate-100 bg-white p-4">
                <div className="flex items-center justify-between">
                  <h3 className="text-sm font-semibold text-slate-900">Internal Team</h3>
                  <div className="flex items-center gap-3">
                    <a className="text-sm text-indigo-600">View all</a>
                    <button className="rounded bg-white px-3 py-1 text-sm border border-slate-200">+ Add Member</button>
                  </div>
                </div>
                <div className="mt-4 overflow-x-auto">
                  <table className="min-w-full text-left text-sm">
                    <thead className="bg-slate-50 text-xs uppercase text-slate-500">
                      <tr>
                        <th className="px-4 py-3">Name</th>
                        <th className="px-4 py-3">Role</th>
                        <th className="px-4 py-3">Department</th>
                        <th className="px-4 py-3">Status</th>
                        <th className="px-4 py-3">Actions</th>
                      </tr>
                    </thead>
                    <tbody>
                      {demo.team.map((m) => (
                        <tr key={m.id} className="border-t border-slate-100">
                          <td className="px-4 py-3">{m.name}</td>
                          <td className="px-4 py-3">{m.role}</td>
                          <td className="px-4 py-3">{m.dept}</td>
                          <td className="px-4 py-3">{m.status==='Active'? <span className="rounded-full bg-emerald-50 px-2 py-1 text-xs text-emerald-700">Active</span>: <span className="rounded-full bg-rose-50 px-2 py-1 text-xs text-rose-700">Inactive</span>}</td>
                          <td className="px-4 py-3">...</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
                <p className="mt-3 text-xs text-slate-500">Total 12 members</p>
              </div>
            </div>
          </section>
        </main>
      </div>
    </div>
  );
}
