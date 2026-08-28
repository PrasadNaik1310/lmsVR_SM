import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { 
  createAcademicSession,
  deleteAcademicSession,
  getAcademicSession,
  listAcademicSessions,
  listBatchesByCourse, 
  createBatchForCourse, 
  getBatchDetails, 
  listCoursesBySession, 
  updateAcademicSession,
  assignCourseToSession,
  listCoursesForUser,
  createCompanyUser,
} from "../api/manageCompany.js";


const demo = {
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
  const [academicSessions, setAcademicSessions] = useState([]);
  const [selectedSessionId, setSelectedSessionId] = useState("");
  const [batches, setBatches] = useState([]);
  const [courses, setCourses] = useState([]);
  const [selectedCourseId, setSelectedCourseId] = useState("");
  const [sessionCourses, setSessionCourses] = useState({});
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [showAcademicSessionModal, setShowAcademicSessionModal] = useState(false);
  const [showCreateBatchModal, setShowCreateBatchModal] = useState(false);
  const [editingAcademicSessionId, setEditingAcademicSessionId] = useState("");
  const [academicSessionForm, setAcademicSessionForm] = useState({
    name: "",
    startDate: "",
    endDate: "",
    isActive: true,
  });
  const [createBatchForm, setCreateBatchForm] = useState({
    courseId: "",
    batchName: "",
    startDate: "",
    endDate: "",
    maxStudents: "",
    status: "draft",
  });
  const token = localStorage.getItem('auth_token');

  const handleLogout = () => {
    localStorage.removeItem('auth_token');
    navigate('/');
  };

  useEffect(() => {
    const fetchAcademicSessions = async () => {
      if (!token) return;
      try {
        setLoading(true);
        const data = await listAcademicSessions(token);
        const sessions = data.academic_sessions || [];
        setAcademicSessions(sessions);
        if (sessions.length > 0) {
          const firstSessionId = sessions[0].id;
          setSelectedSessionId(firstSessionId);
          const sessionCourseData = await listCoursesBySession(firstSessionId, { page: 1, size: 20 }, token);
          setSessionCourses((prev) => ({ ...prev, [firstSessionId]: sessionCourseData.courses || [] }));
        }
      } catch (err) {
        console.error('Failed to fetch academic sessions', err);
      } finally {
        setLoading(false);
      }
    };

    const fetchCourses = async () => {
      if (!token) return;
      try {
        setLoading(true);
        const data = await listCoursesForUser({ page: 1, size: 20 }, token);
        setCourses(data.courses || []);
        // Auto-load batches for first course
        if (data.courses && data.courses.length > 0) {
          setSelectedCourseId(data.courses[0].id);
          setCreateBatchForm((prev) => ({ ...prev, courseId: data.courses[0].id }));
          const courseData = await listBatchesByCourse(data.courses[0].id, { page: 1, size: 20 }, token);
          setBatches(courseData.batches || []);
        }
      } catch (err) {
        console.error('Failed to fetch courses for user', err);
      } finally {
        setLoading(false);
      }
    };

    fetchAcademicSessions();
    fetchCourses();
  }, [token]);

  const openAcademicSessionModal = (session = null) => {
    if (session) {
      setEditingAcademicSessionId(session.id);
      setAcademicSessionForm({
        name: session.name || "",
        startDate: session.start_date ? String(session.start_date).slice(0, 10) : "",
        endDate: session.end_date ? String(session.end_date).slice(0, 10) : "",
        isActive: Boolean(session.is_active),
      });
    } else {
      setEditingAcademicSessionId("");
      setAcademicSessionForm({
        name: "",
        startDate: "",
        endDate: "",
        isActive: true,
      });
    }
    setShowAcademicSessionModal(true);
    setError(null);
  };

  const closeAcademicSessionModal = () => {
    setShowAcademicSessionModal(false);
  };

  const handleAcademicSessionInputChange = (e) => {
    const { name, value, type, checked } = e.target;
    setAcademicSessionForm((prev) => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }));
  };

  const refreshAcademicSessions = async () => {
    const data = await listAcademicSessions(token);
    const sessions = data.academic_sessions || [];
    setAcademicSessions(sessions);
    if (sessions.length === 0) {
      setSelectedSessionId("");
      setSessionCourses({});
      return sessions;
    }
    if (!selectedSessionId || !sessions.some((session) => session.id === selectedSessionId)) {
      setSelectedSessionId(sessions[0].id);
    }
    return sessions;
  };

  const handleSaveAcademicSession = async (e) => {
    e.preventDefault();
    if (!academicSessionForm.name || !academicSessionForm.startDate || !academicSessionForm.endDate) {
      setError('Please fill all academic session fields.');
      return;
    }

    try {
      setLoading(true);
      const payload = {
        name: academicSessionForm.name,
        start_date: academicSessionForm.startDate,
        end_date: academicSessionForm.endDate,
        is_active: academicSessionForm.isActive,
      };

      if (editingAcademicSessionId) {
        await updateAcademicSession(editingAcademicSessionId, payload, token);
      } else {
        await createAcademicSession(payload, token);
      }

      await refreshAcademicSessions();
      closeAcademicSessionModal();
      setError(null);
    } catch (err) {
      setError('Failed to save academic session: ' + err.message);
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleViewAcademicSession = async (sessionId) => {
    try {
      setLoading(true);
      const data = await getAcademicSession(sessionId, token);
      alert('Academic session details: ' + JSON.stringify(data.academic_session, null, 2));
      setError(null);
    } catch (err) {
      setError('Failed to fetch academic session: ' + err.message);
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteAcademicSession = async (sessionId) => {
    const confirmed = window.confirm('Delete this academic session? This does not delete courses.');
    if (!confirmed) return;

    try {
      setLoading(true);
      await deleteAcademicSession(sessionId, token);
      await refreshAcademicSessions();
      if (selectedSessionId === sessionId) {
        setSelectedSessionId("");
        setSessionCourses({});
      }
      setError(null);
    } catch (err) {
      setError('Failed to delete academic session: ' + err.message);
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // Fetch batches by course
  const handleFetchBatchesByCourse = async (courseId) => {
    if (!courseId) {
      setError('Please select a valid course.');
      return;
    }
    try {
      setLoading(true);
      const data = await listBatchesByCourse(courseId, { page: 1, size: 20 }, token);
      setBatches(data.batches || []);
      setSelectedCourseId(courseId);
      setError(null);
    } catch (err) {
      setError('Failed to fetch batches: ' + err.message);
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const openCreateBatchModal = (courseId) => {
    setCreateBatchForm({
      courseId: courseId || selectedCourseId || courses[0]?.id || "",
      batchName: "",
      startDate: "",
      endDate: "",
      maxStudents: "",
      status: "draft",
    });
    setShowCreateBatchModal(true);
    setError(null);
  };

  const closeCreateBatchModal = () => {
    setShowCreateBatchModal(false);
  };

  const handleCreateBatchInputChange = (e) => {
    const { name, value } = e.target;
    setCreateBatchForm((prev) => ({ ...prev, [name]: value }));
  };

  // Create new batch from modal form
  const handleCreateBatch = async (e) => {
    e.preventDefault();
    if (!createBatchForm.courseId || !createBatchForm.batchName || !createBatchForm.startDate || !createBatchForm.endDate || !createBatchForm.maxStudents) {
      setError('Please fill all batch form fields.');
      return;
    }

    try {
      setLoading(true);
      const payload = {
        batch_name: createBatchForm.batchName,
        start_date: createBatchForm.startDate,
        end_date: createBatchForm.endDate,
        max_students: parseInt(createBatchForm.maxStudents, 10),
        status: createBatchForm.status || 'draft',
      };
      await createBatchForCourse(createBatchForm.courseId, payload, token);
      await handleFetchBatchesByCourse(createBatchForm.courseId);
      closeCreateBatchModal();
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
      setSessionCourses((prev) => ({ ...prev, [sessionId]: data.courses || [] }));
      setSelectedSessionId(sessionId);
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
          
          {/*<aside className="hidden w-64 shrink-0 border-r border-slate-100 bg-white lg:block">
            
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
*/}
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
              <button
                type="button"
                onClick={() => openAcademicSessionModal()}
                className="rounded bg-indigo-600 px-4 py-2 text-white"
              >
                + Academic Session
              </button>
              {error && <div className="text-red-600 text-sm">{error}</div>}
              {loading && <div className="text-blue-600 text-sm">Loading...</div>}
            </div>
          </header>

          <section className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4 mb-6">
            <div className="rounded-xl border border-slate-100 bg-white p-4">
              <p className="text-xs text-slate-500">Academic Sessions</p>
              <p className="mt-2 text-2xl font-semibold">{academicSessions.length}</p>
              <p className="text-sm text-slate-400">{academicSessions.filter((session) => session.is_active).length} active</p>
            </div>
            <div className="rounded-xl border border-slate-100 bg-white p-4">
              <p className="text-xs text-slate-500">Courses / Packages</p>
              <p className="mt-2 text-2xl font-semibold">{(courses.length ? courses : demo.courses).length}</p>
              <p className="text-sm text-slate-400">{(courses.length ? courses : demo.courses).filter((course) => (course.status || course.Status || "").toLowerCase() === "draft").length} draft</p>
            </div>
            <div className="rounded-xl border border-slate-100 bg-white p-4">
              <p className="text-xs text-slate-500">Batches</p>
              <p className="mt-2 text-2xl font-semibold">{batches.length}</p>
              <p className="text-sm text-slate-400">{batches.filter((batch) => (batch.status || batch.Status || "").toLowerCase() === "active").length} active</p>
            </div>
            <div className="rounded-xl border border-slate-100 bg-white p-4">
              <p className="text-xs text-slate-500">Team Members</p>
              <p className="mt-2 text-2xl font-semibold">{demo.team.length}</p>
              <p className="text-sm text-slate-400">{demo.team.filter((member) => member.status !== 'Active').length} inactive</p>
            </div>
          </section>

          <section className="grid gap-6 lg:grid-cols-3">
            <div className="lg:col-span-2 space-y-6">
              <div className="rounded-xl border border-slate-100 bg-white p-4">
                <div className="flex items-center justify-between">
                  <h3 className="text-sm font-semibold text-slate-900">Academic Sessions</h3>
                  <div className="flex items-center gap-3">
                    <a
                      onClick={refreshAcademicSessions}
                      className="text-sm text-indigo-600 cursor-pointer"
                    >
                      View all
                    </a>
                    <button
                      type="button"
                      onClick={() => openAcademicSessionModal()}
                      className="rounded bg-white px-3 py-1 text-sm border border-slate-200"
                    >
                      + Add Session
                    </button>
                  </div>
                </div>
                <div className="mt-4 grid gap-3">
                  {academicSessions.length > 0 ? (
                    academicSessions.map((session) => (
                      <div
                        key={session.id}
                        className="rounded-xl border border-slate-100 bg-slate-50 p-3"
                      >
                        <div className="flex items-start justify-between gap-3">
                          <div>
                            <p className="text-xs text-slate-500">
                              {session.is_active ? 'Active' : 'Inactive'}
                            </p>
                            <p className="mt-1 font-semibold text-slate-900">{session.name}</p>
                          </div>

                          <span className={`rounded-full px-2 py-1 text-xs ${session.is_active ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-600'}`}>
                            {session.is_active ? 'Active' : 'Inactive'}
                          </span>
                        </div>

                        <p className="mt-2 text-xs text-slate-500">
                          {session.start_date ? new Date(session.start_date).toLocaleDateString() : ''}
                          {session.start_date && session.end_date ? ' - ' : ''}
                          {session.end_date ? new Date(session.end_date).toLocaleDateString() : ''}
                        </p>

                        <div className="mt-3 flex flex-wrap gap-2">
                          <button
                            type="button"
                            onClick={() => handleListCoursesBySession(session.id)}
                            className="text-sm text-indigo-600 hover:text-indigo-700"
                          >
                            View Courses
                          </button>
                          <button
                            type="button"
                            onClick={() => handleViewAcademicSession(session.id)}
                            className="text-sm text-slate-700 hover:text-slate-900"
                          >
                            View
                          </button>
                          <button
                            type="button"
                            onClick={() => openAcademicSessionModal(session)}
                            className="text-sm text-slate-700 hover:text-slate-900"
                          >
                            Edit
                          </button>
                          <button
                            type="button"
                            onClick={() => handleDeleteAcademicSession(session.id)}
                            className="text-sm text-rose-600 hover:text-rose-700"
                          >
                            Delete
                          </button>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="rounded-xl border border-dashed border-slate-200 bg-slate-50 p-6 text-center text-slate-500">
                      <p className="text-sm">No academic sessions available</p>
                      <button
                        type="button"
                        onClick={() => openAcademicSessionModal()}
                        className="mt-3 rounded bg-indigo-600 px-3 py-1 text-sm text-white"
                      >
                        Create Session
                      </button>
                    </div>
                  )}
                </div>
                <p className="mt-3 text-xs text-slate-500">Total {academicSessions.length} sessions</p>
              </div>

              {selectedSessionId && (
                <div className="rounded-xl border border-slate-100 bg-white p-4">
                  <div className="flex items-center justify-between">
                    <h3 className="text-sm font-semibold text-slate-900">Session Courses</h3>
                    <button
                      type="button"
                      onClick={() => handleListCoursesBySession(selectedSessionId)}
                      className="text-sm text-indigo-600 cursor-pointer"
                    >
                      Refresh
                    </button>
                  </div>

                  <div className="mt-4 grid gap-3">
                    {(sessionCourses[selectedSessionId] || []).length > 0 ? (
                      (sessionCourses[selectedSessionId] || []).map((course) => (
                        <div key={course.id} className="rounded-xl border border-slate-100 bg-slate-50 p-3">
                          <div className="flex items-start justify-between gap-3">
                            <div>
                              <p className="text-xs text-slate-500">{course.level}</p>
                              <p className="mt-1 font-semibold text-slate-900">{course.title}</p>
                            </div>
                            <div className="text-right">
                              <p className="text-sm text-slate-700">
                                {course.booked_seats ? `${course.booked_seats}/${course.total_seats}` : course.total_seats}
                              </p>
                            </div>
                          </div>
                          <p className="mt-2 text-xs text-slate-500">
                            {course.start_date ? new Date(course.start_date).toLocaleDateString() : ''}
                            {course.start_date && course.end_date ? ' - ' : ''}
                            {course.end_date ? new Date(course.end_date).toLocaleDateString() : ''}
                          </p>
                        </div>
                      ))
                    ) : (
                      <div className="rounded-xl border border-dashed border-slate-200 bg-slate-50 p-6 text-center text-slate-500">
                        <p className="text-sm">No courses assigned to this session</p>
                      </div>
                    )}
                  </div>
                  <p className="mt-3 text-xs text-slate-500">
                    Total {(sessionCourses[selectedSessionId] || []).length} courses
                  </p>
                </div>
              )}

              <div className="rounded-xl border border-slate-100 bg-white p-4">
                <div className="flex items-center justify-between">
                  <h3 className="text-sm font-semibold text-slate-900">Batches</h3>
                  <div className="flex items-center gap-3">
                    <a 
                      onClick={() => handleFetchBatchesByCourse(selectedCourseId || courses[0]?.id)}
                      className="text-sm text-indigo-600 cursor-pointer"
                    >
                      View all
                    </a>
                    <button 
                      onClick={() => openCreateBatchModal(selectedCourseId || courses[0]?.id)}
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
                                onClick={() => handleGetBatchDetails(b.course_id || selectedCourseId, b.id || b.ID)}
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

      {showAcademicSessionModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/40 p-4">
          <div className="w-full max-w-lg rounded-xl bg-white p-6 shadow-xl">
            <div className="mb-4 flex items-center justify-between">
              <h3 className="text-lg font-semibold text-slate-900">
                {editingAcademicSessionId ? 'Edit Academic Session' : 'Create Academic Session'}
              </h3>
              <button
                type="button"
                onClick={closeAcademicSessionModal}
                className="rounded px-2 py-1 text-slate-500 hover:bg-slate-100"
              >
                x
              </button>
            </div>

            <form onSubmit={handleSaveAcademicSession} className="space-y-4">
              <div>
                <label className="mb-1 block text-sm font-medium text-slate-700">Session Name</label>
                <input
                  type="text"
                  name="name"
                  value={academicSessionForm.name}
                  onChange={handleAcademicSessionInputChange}
                  className="w-full rounded border border-slate-300 px-3 py-2 text-sm outline-none focus:border-indigo-500"
                  placeholder="e.g. 2026 - 2027"
                  required
                />
              </div>

              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <div>
                  <label className="mb-1 block text-sm font-medium text-slate-700">Start Date</label>
                  <input
                    type="date"
                    name="startDate"
                    value={academicSessionForm.startDate}
                    onChange={handleAcademicSessionInputChange}
                    className="w-full rounded border border-slate-300 px-3 py-2 text-sm outline-none focus:border-indigo-500"
                    required
                  />
                </div>
                <div>
                  <label className="mb-1 block text-sm font-medium text-slate-700">End Date</label>
                  <input
                    type="date"
                    name="endDate"
                    value={academicSessionForm.endDate}
                    onChange={handleAcademicSessionInputChange}
                    className="w-full rounded border border-slate-300 px-3 py-2 text-sm outline-none focus:border-indigo-500"
                    required
                  />
                </div>
              </div>

              <div>
                <label className="mb-1 block text-sm font-medium text-slate-700">Status</label>
                <select
                  name="isActive"
                  value={academicSessionForm.isActive ? 'true' : 'false'}
                  onChange={(e) =>
                    setAcademicSessionForm((prev) => ({
                      ...prev,
                      isActive: e.target.value === 'true',
                    }))
                  }
                  className="w-full rounded border border-slate-300 px-3 py-2 text-sm outline-none focus:border-indigo-500"
                >
                  <option value="true">Active</option>
                  <option value="false">Inactive</option>
                </select>
              </div>

              <div className="flex items-center justify-end gap-3 pt-2">
                <button
                  type="button"
                  onClick={closeAcademicSessionModal}
                  className="rounded border border-slate-300 px-4 py-2 text-sm text-slate-700 hover:bg-slate-50"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="rounded bg-indigo-600 px-4 py-2 text-sm text-white hover:bg-indigo-700 disabled:opacity-50"
                  disabled={loading}
                >
                  {loading ? 'Saving...' : (editingAcademicSessionId ? 'Update Session' : 'Create Session')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showCreateBatchModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/40 p-4">
          <div className="w-full max-w-lg rounded-xl bg-white p-6 shadow-xl">
            <div className="mb-4 flex items-center justify-between">
              <h3 className="text-lg font-semibold text-slate-900">Create Batch</h3>
              <button
                type="button"
                onClick={closeCreateBatchModal}
                className="rounded px-2 py-1 text-slate-500 hover:bg-slate-100"
              >
                x
              </button>
            </div>

            <form onSubmit={handleCreateBatch} className="space-y-4">
              <div>
                <label className="mb-1 block text-sm font-medium text-slate-700">Course</label>
                <select
                  name="courseId"
                  value={createBatchForm.courseId}
                  onChange={handleCreateBatchInputChange}
                  className="w-full rounded border border-slate-300 px-3 py-2 text-sm outline-none focus:border-indigo-500"
                  required
                >
                  <option value="">Select course</option>
                  {courses.map((course, index) => (
                    <option key={course.id || course.ID || `course-option-${index}`} value={course.id || course.ID}>
                      {course.title || course.Title || `Course ${index + 1}`}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="mb-1 block text-sm font-medium text-slate-700">Batch Name</label>
                <input
                  type="text"
                  name="batchName"
                  value={createBatchForm.batchName}
                  onChange={handleCreateBatchInputChange}
                  className="w-full rounded border border-slate-300 px-3 py-2 text-sm outline-none focus:border-indigo-500"
                  placeholder="e.g. Batch - Fall 2026"
                  required
                />
              </div>

              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <div>
                  <label className="mb-1 block text-sm font-medium text-slate-700">Start Date</label>
                  <input
                    type="date"
                    name="startDate"
                    value={createBatchForm.startDate}
                    onChange={handleCreateBatchInputChange}
                    className="w-full rounded border border-slate-300 px-3 py-2 text-sm outline-none focus:border-indigo-500"
                    required
                  />
                </div>
                <div>
                  <label className="mb-1 block text-sm font-medium text-slate-700">End Date</label>
                  <input
                    type="date"
                    name="endDate"
                    value={createBatchForm.endDate}
                    onChange={handleCreateBatchInputChange}
                    className="w-full rounded border border-slate-300 px-3 py-2 text-sm outline-none focus:border-indigo-500"
                    required
                  />
                </div>
              </div>

              <div>
                <label className="mb-1 block text-sm font-medium text-slate-700">Max Students</label>
                <input
                  type="number"
                  name="maxStudents"
                  min="1"
                  value={createBatchForm.maxStudents}
                  onChange={handleCreateBatchInputChange}
                  className="w-full rounded border border-slate-300 px-3 py-2 text-sm outline-none focus:border-indigo-500"
                  placeholder="e.g. 40"
                  required
                />
              </div>

              <div>
                <label className="mb-1 block text-sm font-medium text-slate-700">Status</label>
                <select
                  name="status"
                  value={createBatchForm.status}
                  onChange={handleCreateBatchInputChange}
                  className="w-full rounded border border-slate-300 px-3 py-2 text-sm outline-none focus:border-indigo-500"
                >
                  <option value="draft">Draft</option>
                  <option value="active">Active</option>
                  <option value="inactive">Inactive</option>
                </select>
              </div>

              <div className="flex items-center justify-end gap-3 pt-2">
                <button
                  type="button"
                  onClick={closeCreateBatchModal}
                  className="rounded border border-slate-300 px-4 py-2 text-sm text-slate-700 hover:bg-slate-50"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="rounded bg-indigo-600 px-4 py-2 text-sm text-white hover:bg-indigo-700 disabled:opacity-50"
                  disabled={loading || courses.length === 0}
                >
                  {loading ? 'Creating...' : 'Create Batch'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
  const [showAddMemberModal, setShowAddMemberModal] = useState(false);

const [memberForm, setMemberForm] = useState({
  first_name: "",
  last_name: "",
  email: "",
  phone: "",
  password: "",
  role: "teacher",

  // student
  enrollment_number: "",
  date_of_birth: "",
  address: "",
  admission_date: "",

  // teacher
  specialization: "",
  bio: "",
  joining_date: "",
});
}

