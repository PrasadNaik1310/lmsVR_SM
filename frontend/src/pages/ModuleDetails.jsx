// src/pages/ModuleDetails.jsx

import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { http } from "../services/http";

export default function ModuleDetails() {
  const { id } = useParams();

  const token = localStorage.getItem("auth_token");

  const [lessons, setLessons] = useState([]);

  const [form, setForm] = useState({
    title: "",
    description: "",
    content_type: "PDF",
    duration_minutes: 30,
  });

  async function fetchLessons() {
    try {
      const res = await http.get(
        `/modules/${id}/lessons`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      setLessons(res.data.lessons || []);
    } catch (err) {
      console.error(err);
      alert("Failed to fetch lessons");
    }
  }

  async function createLesson(e) {
    e.preventDefault();

    try {
      await http.post(
        `/modules/${id}/lessons`,
        {
          title: form.title,
          description: form.description,
          content_type: form.content_type,
          duration_minutes: Number(
            form.duration_minutes
          ),
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      alert("Lesson created");

      setForm({
        title: "",
        description: "",
        content_type: "PDF",
        duration_minutes: 30,
      });

      fetchLessons();
    } catch (err) {
      console.error(err);

      alert(
        err?.response?.data?.error ||
          "Failed creating lesson"
      );
    }
  }

  async function deleteLesson(lessonID) {
    try {
      await http.delete(
        `/lessons/${lessonID}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      alert("Lesson deleted");

      fetchLessons();
    } catch (err) {
      console.error(err);

      alert(
        err?.response?.data?.error ||
          "Failed deleting lesson"
      );
    }
  }

  useEffect(() => {
    fetchLessons();
  }, []);

  return (
    <div className="p-6 space-y-8">

      <div className="bg-white shadow rounded p-6">
        <h1 className="text-3xl font-bold">
          Module Lessons
        </h1>
      </div>

      <div className="bg-white shadow rounded p-6">

        <h2 className="text-xl font-semibold mb-4">
          Create Lesson
        </h2>

        <form
          onSubmit={createLesson}
          className="space-y-4"
        >
          <input
            type="text"
            placeholder="Lesson Title"
            value={form.title}
            onChange={(e) =>
              setForm({
                ...form,
                title: e.target.value,
              })
            }
            className="w-full border rounded p-2"
            required
          />

          <textarea
            placeholder="Description"
            value={form.description}
            onChange={(e) =>
              setForm({
                ...form,
                description: e.target.value,
              })
            }
            className="w-full border rounded p-2"
            rows={4}
          />

          <select
            value={form.content_type}
            onChange={(e) =>
              setForm({
                ...form,
                content_type: e.target.value,
              })
            }
            className="w-full border rounded p-2"
          >
            <option value="PDF">PDF</option>
          </select>

          <input
            type="number"
            min="1"
            placeholder="Duration Minutes"
            value={form.duration_minutes}
            onChange={(e) =>
              setForm({
                ...form,
                duration_minutes:
                  e.target.value,
              })
            }
            className="w-full border rounded p-2"
          />

          <button
            type="submit"
            className="bg-green-600 text-white px-4 py-2 rounded"
          >
            Create Lesson
          </button>

        </form>
      </div>

      <div className="bg-white shadow rounded p-6">

        <h2 className="text-xl font-semibold mb-4">
          Lessons
        </h2>

        <table className="w-full">
          <thead>
            <tr className="border-b">
              <th className="text-left p-2">
                Title
              </th>

              <th className="text-left p-2">
                Type
              </th>

              <th className="text-left p-2">
                Duration
              </th>

              <th className="text-left p-2">
                Published
              </th>

              <th className="text-left p-2">
                Actions
              </th>
            </tr>
          </thead>

          <tbody>
            {lessons.map((lesson) => (
              <tr
                key={lesson.id}
                className="border-b"
              >
                <td className="p-2">
                  {lesson.title}
                </td>

                <td className="p-2">
                  {lesson.content_type}
                </td>

                <td className="p-2">
                  {lesson.duration_minutes}
                </td>

                <td className="p-2">
                  {lesson.is_published
                    ? "Yes"
                    : "No"}
                </td>

                <td className="p-2">
                  <button
                    onClick={() =>
                      deleteLesson(lesson.id)
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

      </div>

    </div>
  );
}