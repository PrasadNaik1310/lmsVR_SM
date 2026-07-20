import { useEffect, useState } from "react";
import { http } from "../services/http";

export default function CourseSchedule({ courseId }) {
  const token = localStorage.getItem("auth_token");

  const [schedules, setSchedules] = useState([]);
  const [modules, setModules] = useState([]);
  const [lessons, setLessons] = useState([]);
  const [teachers, setTeachers] = useState([]);

  const [form, setForm] = useState({
    lesson_id: "",
    teacher_id: "",
    planned_date: "",
    start_time: "",
    end_time: "",
  });

  async function fetchSchedules() {
    try {
      const res = await http.get(`/courses/${courseId}/schedules`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      setSchedules(res.data.schedules || []);
    } catch (err) {
      console.error(err);
      alert("Failed to fetch schedules");
    }
  }

  async function fetchModules() {
    try {
      const res = await http.get(`/courses/${courseId}/modules`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      const modulesData = res.data.modules || [];
      setModules(modulesData);

      let allLessons = [];

      for (const module of modulesData) {
        try {
          const lessonRes = await http.get(
            `/modules/${module.id}/lessons`,
            {
              headers: {
                Authorization: `Bearer ${token}`,
              },
            }
          );

          allLessons = [
            ...allLessons,
            ...(lessonRes.data.lessons || []),
          ];
        } catch (err) {
          console.error(
            `Failed loading lessons for module ${module.id}`,
            err
          );
        }
      }

      setLessons(allLessons);
    } catch (err) {
      console.error(err);
      alert("Failed to fetch modules");
    }
  }

  async function fetchTeachers() {
    try {
      const res = await http.get("/teachers", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      setTeachers(res.data.teachers || []);
    } catch (err) {
      console.error(err);
      alert("Failed to fetch teachers");
    }
  }

  async function handleCreateSchedule(e) {
    e.preventDefault();

    try {
      await http.post(
        `/courses/${courseId}/schedules`,
        {
          lesson_id: form.lesson_id,
          teacher_id: form.teacher_id,
          planned_date: form.planned_date,
          start_time: form.start_time,
          end_time: form.end_time,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      alert("Schedule created");

      setForm({
        lesson_id: "",
        teacher_id: "",
        planned_date: "",
        start_time: "",
        end_time: "",
      });

      fetchSchedules();
    } catch (err) {
      console.error(err);

      alert(
        err?.response?.data?.error ||
          "Failed to create schedule"
      );
    }
  }

  useEffect(() => {
    fetchSchedules();
    fetchModules();
    fetchTeachers();
  }, [courseId]);

  return (
    <div className="space-y-8">
      <div className="bg-white shadow rounded p-6">
        <h2 className="text-xl font-semibold mb-4">
          Create Schedule
        </h2>

        <form
          onSubmit={handleCreateSchedule}
          className="space-y-4"
        >
          <select
            value={form.lesson_id}
            onChange={(e) =>
              setForm({
                ...form,
                lesson_id: e.target.value,
              })
            }
            className="w-full border rounded p-2"
            required
          >
            <option value="">Select Lesson</option>

            {lessons.map((lesson) => (
              <option
                key={lesson.id}
                value={lesson.id}
              >
                {lesson.title}
              </option>
            ))}
          </select>

          <select
            value={form.teacher_id}
            onChange={(e) =>
              setForm({
                ...form,
                teacher_id: e.target.value,
              })
            }
            className="w-full border rounded p-2"
            required
          >
            <option value="">Select Teacher</option>

            {teachers.map((teacher) => (
              <option
                key={teacher.id}
                value={teacher.id}
              >
                {teacher.name}
              </option>
            ))}
          </select>

          <input
            type="date"
            value={form.planned_date}
            onChange={(e) =>
              setForm({
                ...form,
                planned_date: e.target.value,
              })
            }
            className="w-full border rounded p-2"
            required
          />

          <input
            type="time"
            value={form.start_time}
            onChange={(e) =>
              setForm({
                ...form,
                start_time: e.target.value,
              })
            }
            className="w-full border rounded p-2"
            required
          />

          <input
            type="time"
            value={form.end_time}
            onChange={(e) =>
              setForm({
                ...form,
                end_time: e.target.value,
              })
            }
            className="w-full border rounded p-2"
            required
          />

          <button
            type="submit"
            className="bg-green-600 text-white px-4 py-2 rounded"
          >
            Create Schedule
          </button>
        </form>
      </div>

      <div className="bg-white shadow rounded p-6">
        <h2 className="text-xl font-semibold mb-4">
          Scheduled Sessions
        </h2>

        <table className="w-full">
          <thead>
            <tr className="border-b">
              <th className="text-left p-2">Lesson</th>
              <th className="text-left p-2">Teacher</th>
              <th className="text-left p-2">Date</th>
              <th className="text-left p-2">Start</th>
              <th className="text-left p-2">End</th>
              <th className="text-left p-2">Status</th>
            </tr>
          </thead>

          <tbody>
            {schedules.map((schedule) => (
              <tr
                key={schedule.id}
                className="border-b"
              >
                <td className="p-2">
                  {schedule.lesson_title}
                </td>

                <td className="p-2">
                  {schedule.teacher_name}
                </td>

                <td className="p-2">
                  {new Date(
                    schedule.planned_date
                  ).toLocaleDateString()}
                </td>

                <td className="p-2">
                  {schedule.planned_start_time}
                </td>

                <td className="p-2">
                  {schedule.planned_end_time}
                </td>

                <td className="p-2">
                  {schedule.status}
                </td>
              </tr>
            ))}

            {schedules.length === 0 && (
              <tr>
                <td
                  colSpan="6"
                  className="text-center p-4 text-gray-500"
                >
                  No schedules found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}