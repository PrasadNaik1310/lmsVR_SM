import { useEffect, useState } from "react";
import {
  listAcademicSessions,
  createAcademicSession,
  updateAcademicSession,
  deleteAcademicSession,
  listCoursesBySession,
  assignCourseToSession,
  listCoursesForUser,
} from "../api/manageCompany";

export default function AcademicSessionManagement() {
  const token = localStorage.getItem("auth_token");

  const [sessions, setSessions] = useState([]);
  const [courses, setCourses] = useState([]);
  const [sessionCourses, setSessionCourses] = useState([]);

  const [loading, setLoading] = useState(false);
  const [coursesLoading, setCoursesLoading] = useState(false);

  const [selectedSession, setSelectedSession] = useState(null);
  const [editingSession, setEditingSession] = useState(null);

  const [form, setForm] = useState({
    name: "",
    start_date: "",
    end_date: "",
    is_active: true,
  });

  async function fetchSessions() {
    try {
      setLoading(true);

      const data = await listAcademicSessions(token);

      setSessions(data.sessions || []);
    } catch (err) {
      console.error("Failed to load academic sessions:", err);
      console.error("Status:", err.response?.status);
      console.error("Data:", err.response?.data);

      alert(
        err?.response?.data?.error ||
          "Failed to load academic sessions"
      );
    } finally {
      setLoading(false);
    }
  }

  async function fetchCourses() {
    try {
      const data = await listCoursesForUser({}, token);

      setCourses(data.courses || []);
    } catch (err) {
      console.error("Failed to load courses:", err);

      alert(
        err?.response?.data?.error ||
          "Failed to load courses"
      );
    }
  }

  async function fetchSessionCourses(sessionId) {
    try {
      setCoursesLoading(true);

      const data = await listCoursesBySession(
        sessionId,
        {},
        token
      );

      setSessionCourses(data.courses || []);
    } catch (err) {
      console.error("Failed to load session courses:", err);

      alert(
        err?.response?.data?.error ||
          "Failed to load session courses"
      );
    } finally {
      setCoursesLoading(false);
    }
  }

  useEffect(() => {
    fetchSessions();
    fetchCourses();
  }, []);

  function handleChange(e) {
    const { name, value, type, checked } = e.target;

    setForm((prev) => ({
      ...prev,
      [name]: type === "checkbox" ? checked : value,
    }));
  }

  function resetForm() {
    setForm({
      name: "",
      start_date: "",
      end_date: "",
      is_active: true,
    });

    setEditingSession(null);
  }

  async function handleCreateSession(e) {
    e.preventDefault();

    try {
      await createAcademicSession(
        {
          name: form.name,
          start_date: form.start_date,
          end_date: form.end_date,
          is_active: form.is_active,
        },
        token
      );

      alert("Academic session created successfully");

      resetForm();
      fetchSessions();
    } catch (err) {
      console.error(err);

      alert(
        err?.response?.data?.error ||
          "Failed to create academic session"
      );
    }
  }

  function handleEdit(session) {
    setEditingSession(session);

    setForm({
      name: session.name || "",
      start_date: session.start_date
        ? session.start_date.substring(0, 10)
        : "",
      end_date: session.end_date
        ? session.end_date.substring(0, 10)
        : "",
      is_active: session.is_active ?? true,
    });

    window.scrollTo({
      top: 0,
      behavior: "smooth",
    });
  }

  async function handleUpdateSession(e) {
    e.preventDefault();

    if (!editingSession) {
      return;
    }

    try {
      await updateAcademicSession(
        editingSession.id,
        {
          name: form.name,
          start_date: form.start_date,
          end_date: form.end_date,
          is_active: form.is_active,
        },
        token
      );

      alert("Academic session updated successfully");

      resetForm();
      fetchSessions();
    } catch (err) {
      console.error(err);

      alert(
        err?.response?.data?.error ||
          "Failed to update academic session"
      );
    }
  }

  async function handleDeleteSession(sessionId) {
    const confirmed = window.confirm(
      "Are you sure you want to delete this academic session?"
    );

    if (!confirmed) {
      return;
    }

    try {
      await deleteAcademicSession(sessionId, token);

      alert("Academic session deleted successfully");

      if (selectedSession?.id === sessionId) {
        setSelectedSession(null);
        setSessionCourses([]);
      }

      fetchSessions();
    } catch (err) {
      console.error(err);

      alert(
        err?.response?.data?.error ||
          "Failed to delete academic session"
      );
    }
  }

  async function handleViewCourses(session) {
    setSelectedSession(session);

    await fetchSessionCourses(session.id);
  }

  async function handleAssignCourse(courseId) {
    if (!selectedSession) {
      return;
    }

    try {
      await assignCourseToSession(
        selectedSession.id,
        courseId,
        token
      );

      alert("Course assigned successfully");

      fetchSessionCourses(selectedSession.id);
    } catch (err) {
      console.error(err);

      alert(
        err?.response?.data?.error ||
          "Failed to assign course"
      );
    }
  }

  return (
    <div className="p-6 space-y-8">

      <div>
        <h1 className="text-3xl font-bold">
          Academic Session Management
        </h1>
      </div>

      {/* CREATE / EDIT SESSION */}

      <div className="bg-white rounded-lg shadow p-6">

        <h2 className="text-xl font-semibold mb-4">
          {editingSession
            ? "Edit Academic Session"
            : "Create Academic Session"}
        </h2>

        <form
          onSubmit={
            editingSession
              ? handleUpdateSession
              : handleCreateSession
          }
          className="grid grid-cols-2 gap-4"
        >

          <input
            name="name"
            placeholder="Academic Session Name"
            value={form.name}
            onChange={handleChange}
            required
            className="border p-2 rounded"
          />

          <div />

          <div>
            <label className="block text-sm mb-1">
              Start Date
            </label>

            <input
              type="date"
              name="start_date"
              value={form.start_date}
              onChange={handleChange}
              required
              className="border p-2 rounded w-full"
            />
          </div>

          <div>
            <label className="block text-sm mb-1">
              End Date
            </label>

            <input
              type="date"
              name="end_date"
              value={form.end_date}
              onChange={handleChange}
              required
              className="border p-2 rounded w-full"
            />
          </div>

          <label className="flex items-center gap-2">
            <input
              type="checkbox"
              name="is_active"
              checked={form.is_active}
              onChange={handleChange}
            />

            Active
          </label>

          <div className="flex gap-2">

            <button
              type="submit"
              className="bg-blue-600 text-white px-4 py-2 rounded"
            >
              {editingSession
                ? "Update Session"
                : "Create Session"}
            </button>

            {editingSession && (
              <button
                type="button"
                onClick={resetForm}
                className="border px-4 py-2 rounded"
              >
                Cancel
              </button>
            )}

          </div>

        </form>
      </div>

      {/* SESSION LIST */}

      <div className="bg-white rounded-lg shadow p-6">

        <div className="flex justify-between items-center mb-4">

          <h2 className="text-xl font-semibold">
            Academic Sessions
          </h2>

          <button
            onClick={fetchSessions}
            className="border px-4 py-2 rounded"
          >
            Refresh
          </button>

        </div>

        {loading ? (
          <p>Loading academic sessions...</p>
        ) : sessions.length === 0 ? (
          <p>No academic sessions found.</p>
        ) : (
          <table className="w-full border-collapse">

            <thead>
              <tr className="border-b">

                <th className="text-left p-3">
                  Name
                </th>

                <th className="text-left p-3">
                  Start Date
                </th>

                <th className="text-left p-3">
                  End Date
                </th>

                <th className="text-left p-3">
                  Status
                </th>

                <th className="text-left p-3">
                  Actions
                </th>

              </tr>
            </thead>

            <tbody>

              {sessions.map((session) => (

                <tr
                  key={session.id}
                  className="border-b"
                >

                  <td className="p-3">
                    {session.name}
                  </td>

                  <td className="p-3">
                    {session.start_date
                      ? session.start_date.substring(0, 10)
                      : "-"}
                  </td>

                  <td className="p-3">
                    {session.end_date
                      ? session.end_date.substring(0, 10)
                      : "-"}
                  </td>

                  <td className="p-3">
                    {session.is_active
                      ? "Active"
                      : "Inactive"}
                  </td>

                  <td className="p-3 flex gap-2">

                    <button
                      onClick={() =>
                        handleViewCourses(session)
                      }
                      className="bg-blue-600 text-white px-3 py-1 rounded"
                    >
                      Courses
                    </button>

                    <button
                      onClick={() =>
                        handleEdit(session)
                      }
                      className="bg-yellow-500 text-white px-3 py-1 rounded"
                    >
                      Edit
                    </button>

                    <button
                      onClick={() =>
                        handleDeleteSession(session.id)
                      }
                      className="bg-red-600 text-white px-3 py-1 rounded"
                    >
                      Delete
                    </button>

                  </td>

                </tr>

              ))}

            </tbody>

          </table>
        )}

      </div>

      {/* SESSION COURSES */}

      {selectedSession && (

        <div className="bg-white rounded-lg shadow p-6">

          <div className="flex justify-between items-center mb-4">

            <h2 className="text-xl font-semibold">
              Courses in {selectedSession.name}
            </h2>

            <button
              onClick={() =>
                fetchSessionCourses(selectedSession.id)
              }
              className="border px-4 py-2 rounded"
            >
              Refresh
            </button>

          </div>

          {coursesLoading ? (
            <p>Loading courses...</p>
          ) : (
            <>
              <div className="mb-6">

                <h3 className="font-semibold mb-3">
                  Assign Course
                </h3>

                <div className="flex gap-2">

                  <select
                    id="course-select"
                    className="border p-2 rounded flex-1"
                    defaultValue=""
                  >

                    <option value="" disabled>
                      Select a course
                    </option>

                    {courses
                      .filter(
                        (course) =>
                          !sessionCourses.some(
                            (sessionCourse) =>
                              sessionCourse.id ===
                              course.id
                          )
                      )
                      .map((course) => (

                        <option
                          key={course.id}
                          value={course.id}
                        >
                          {course.title}
                        </option>

                      ))}

                  </select>

                  <button
                    onClick={() => {
                      const select =
                        document.getElementById(
                          "course-select"
                        );

                      if (!select.value) {
                        alert("Select a course first");
                        return;
                      }

                      handleAssignCourse(
                        select.value
                      );
                    }}
                    className="bg-green-600 text-white px-4 py-2 rounded"
                  >
                    Assign Course
                  </button>

                </div>

              </div>

              <h3 className="font-semibold mb-3">
                Assigned Courses
              </h3>

              {sessionCourses.length === 0 ? (
                <p>
                  No courses assigned to this session.
                </p>
              ) : (
                <table className="w-full border-collapse">

                  <thead>
                    <tr className="border-b">

                      <th className="text-left p-3">
                        Title
                      </th>

                      <th className="text-left p-3">
                        Level
                      </th>

                      <th className="text-left p-3">
                        Status
                      </th>

                    </tr>
                  </thead>

                  <tbody>

                    {sessionCourses.map(
                      (course) => (

                        <tr
                          key={course.id}
                          className="border-b"
                        >

                          <td className="p-3">
                            {course.title}
                          </td>

                          <td className="p-3">
                            {course.level}
                          </td>

                          <td className="p-3">
                            {course.status}
                          </td>

                        </tr>

                      )
                    )}

                  </tbody>

                </table>
              )}

            </>
          )}

        </div>

      )}

    </div>
  );
}