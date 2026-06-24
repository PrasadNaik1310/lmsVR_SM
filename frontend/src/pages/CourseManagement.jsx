// src/pages/CourseManagement.jsx

import { useEffect, useState } from "react";
import {
  createCourse,
  listCourses,
  publishCourse,
  generateInvite,
} from "../api/course";
import { useNavigate } from "react-router-dom"
export default function CourseManagement() {
  const token = localStorage.getItem("auth_token");

  const [courses, setCourses] = useState([]);
  const [loading, setLoading] = useState(false);

  const [form, setForm] = useState({
    title: "",
    description: "",
    level: "",
    total_seats: "",
    meet_link: "",
    start_date: "",
    end_date: "",
  });

  const navigate = useNavigate();

  async function fetchCourses() {
    try {
      setLoading(true);

      const data = await listCourses({}, token);

      setCourses(data.courses || []);
    } catch (err) {
      console.error("Status", err.response?.status);
      console.log("DATA,", err.response?.data)
      alert("Failed to load courses");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    fetchCourses();
  }, []);

  function handleChange(e) {
    setForm((prev) => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  }

  async function handleCreateCourse(e) {
    e.preventDefault();

    try {
      await createCourse(
        {
          title: form.title,
          description: form.description,
          level: form.level,
          total_seats: Number(form.total_seats),
          meet_link: form.meet_link,
          start_date: form.start_date,
          end_date: form.end_date,
        },
        token
      );

      alert("Course created successfully");

      setForm({
        title: "",
        description: "",
        level: "",
        total_seats: "",
        meet_link: "",
        start_date: "",
        end_date: "",
      });

      fetchCourses();
    } catch (err) {
      console.error(err);
      alert(
        err?.response?.data?.error ||
          "Failed to create course"
      );
    }
  }

  async function handlePublish(courseId) {
    try {
      await publishCourse(courseId, token);

      alert("Course published");

      fetchCourses();
    } catch (err) {
      console.error(err);

      alert(
        err?.response?.data?.error ||
          "Failed to publish course"
      );
    }
  }

  async function handleGenerateInvite(courseId) {
    try {
      const data = await generateInvite(courseId, token);

      window.prompt(
        "Copy Invite Link",
        data.invite_link
      );
    } catch (err) {
      console.error(err);

      alert(
        err?.response?.data?.error ||
          "Failed generating invite"
      );
    }
  }

  return (
    <div className="p-6 space-y-8">

      <div>
        <h1 className="text-3xl font-bold">
          Course Management
        </h1>
      </div>

      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-xl font-semibold mb-4">
          Create Course
        </h2>

        <form
          onSubmit={handleCreateCourse}
          className="grid grid-cols-2 gap-4"
        >
          <input
            name="title"
            placeholder="Course Title"
            value={form.title}
            onChange={handleChange}
            required
            className="border p-2 rounded"
          />

          <input
            name="level"
            placeholder="Level"
            value={form.level}
            onChange={handleChange}
            className="border p-2 rounded"
          />

          <input
            name="total_seats"
            type="number"
            placeholder="Total Seats"
            value={form.total_seats}
            onChange={handleChange}
            className="border p-2 rounded"
          />

          <input
            name="meet_link"
            placeholder="Google Meet Link"
            value={form.meet_link}
            onChange={handleChange}
            className="border p-2 rounded"
          />

          <input
            type="date"
            name="start_date"
            value={form.start_date}
            onChange={handleChange}
            className="border p-2 rounded"
          />

          <input
            type="date"
            name="end_date"
            value={form.end_date}
            onChange={handleChange}
            className="border p-2 rounded"
          />

          <textarea
            name="description"
            placeholder="Description"
            value={form.description}
            onChange={handleChange}
            className="border p-2 rounded col-span-2"
            rows={4}
          />

          <button
            type="submit"
            className="bg-blue-600 text-white px-4 py-2 rounded"
          >
            Create Course
          </button>
        </form>
      </div>

      <div className="bg-white rounded-lg shadow p-6">

        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold">
            Courses
          </h2>

          <button
            onClick={fetchCourses}
            className="border px-4 py-2 rounded"
          >
            Refresh
          </button>
        </div>

        {loading ? (
          <p>Loading courses...</p>
        ) : (
          <table className="w-full border-collapse">
            <thead>
              <tr className="border-b">
                <th className="text-left p-3">Title</th>
                <th className="text-left p-3">Level</th>
                <th className="text-left p-3">Status</th>
                <th className="text-left p-3">Seats</th>
                <th className="text-left p-3">Actions</th>
              </tr>
            </thead>

            <tbody>
              {courses.map((course) => (
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

                  <td className="p-3">
                    {course.booked_seats}/
                    {course.total_seats}
                  </td>

                  <td className="p-3 flex gap-2">
                    <button
                      onClick={() =>
                        handlePublish(course.id)
                      }
                      className="bg-green-600 text-white px-3 py-1 rounded"
                    >
                      Publish
                    </button>

                    <button
                      onClick={() =>
                        handleGenerateInvite(course.id)
                      }
                      className="bg-purple-600 text-white px-3 py-1 rounded"
                    >
                      Invite
                    </button>
                    <button
  onClick={() => navigate(`/courses/${course.id}`)}
  className="bg-blue-600 text-white px-3 py-1 rounded"
>
  Modules
</button>
{/*<button
  onClick={() =>
    navigate(`/courses/${course.id}`)}
     className="bg-blue-600 text-white px-3 py-1 rounded"
>
  Lessons
</button>*/}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

    </div>
  );
}