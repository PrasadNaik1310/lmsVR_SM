import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { http } from "../services/http";

export default function CourseDetails() {
  const { id } = useParams();
  const navigate = useNavigate();

  const token = localStorage.getItem("auth_token");

  const [course, setCourse] = useState(null);
  const [modules, setModules] = useState([]);

  const [form, setForm] = useState({
    title: "",
    description: "",
    position: 1,
  });

  async function fetchCourse() {
    try {
      const res = await http.get(`/courses/${id}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      setCourse(res.data.course || res.data);
    } catch (err) {
      console.error(err);
      alert("Failed to fetch course");
    }
  }

  async function fetchModules() {
    try {
      const res = await http.get(`/courses/${id}/modules`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      setModules(res.data.modules || []);
    } catch (err) {
      console.error(err);
      alert("Failed to fetch modules");
    }
  }

  async function handleCreateModule(e) {
    e.preventDefault();

    try {
      await http.post(
        `/courses/${id}/modules`,
        {
          title: form.title,
          description: form.description,
          position: Number(form.position),
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      alert("Module created");

      setForm({
        title: "",
        description: "",
        position: modules.length + 2,
      });

      fetchModules();
    } catch (err) {
      console.error(err);

      alert(
        err?.response?.data?.error ||
          "Failed creating module"
      );
    }
  }

  useEffect(() => {
    fetchCourse();
    fetchModules();
  }, []);

  return (
    <div className="p-6 space-y-8">

      {/* Course Details */}

      <div className="bg-white shadow rounded p-6">
        <h1 className="text-3xl font-bold">
          {course?.title}
        </h1>

        <p className="mt-2 text-gray-600">
          {course?.description}
        </p>

        <div className="mt-4">
          <span className="font-semibold">
            Status:
          </span>{" "}
          {course?.status}
        </div>

        <div className="mt-2">
          <span className="font-semibold">
            Level:
          </span>{" "}
          {course?.level}
        </div>
      </div>

      {/* Create Module */}

      <div className="bg-white shadow rounded p-6">
        <h2 className="text-xl font-semibold mb-4">
          Create Module
        </h2>

        <form
          onSubmit={handleCreateModule}
          className="space-y-4"
        >
          <input
            type="text"
            placeholder="Module Title"
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

          <input
            type="number"
            min="1"
            placeholder="Position"
            value={form.position}
            onChange={(e) =>
              setForm({
                ...form,
                position: e.target.value,
              })
            }
            className="w-full border rounded p-2"
          />

          <button
            type="submit"
            className="bg-green-600 text-white px-4 py-2 rounded"
          >
            Add Module
          </button>
        </form>
      </div>

      {/* Modules List */}

      <div className="bg-white shadow rounded p-6">
        <h2 className="text-xl font-semibold mb-4">
          Modules
        </h2>

        <table className="w-full">
          <thead>
            <tr className="border-b">
              <th className="text-left p-2">
                Position
              </th>

              <th className="text-left p-2">
                Title
              </th>

              <th className="text-left p-2">
                Description
              </th>

              <th className="text-left p-2">
                Actions
              </th>
            </tr>
          </thead>

          <tbody>
            {modules.map((module) => (
              <tr
                key={module.id}
                className="border-b"
              >
                <td className="p-2">
                  {module.position}
                </td>

                <td className="p-2">
                  {module.title}
                </td>

                <td className="p-2">
                  {module.description}
                </td>

                <td className="p-2">
                  <button
                    onClick={() =>
                      navigate(`/modules/${module.id}`)
                    }
                    className="bg-blue-600 text-white px-3 py-1 rounded"
                  >
                    Lessons
                  </button>
                </td>
              </tr>
            ))}

            {modules.length === 0 && (
              <tr>
                <td
                  colSpan="4"
                  className="text-center p-4 text-gray-500"
                >
                  No modules found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

    </div>
  );
}