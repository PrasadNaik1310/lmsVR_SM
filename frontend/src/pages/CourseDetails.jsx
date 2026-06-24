import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { http } from "../services/http";
import { useNavigate } from "react-router-dom";
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

      setCourse(res.data.lesson || res.data.course || res.data);
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
    </tbody>
  </table>

</div>

    </div>
  );
}