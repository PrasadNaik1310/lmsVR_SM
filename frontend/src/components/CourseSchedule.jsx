import { useEffect, useState } from "react";
import { http } from "../services/http";

export default function CourseSchedule({ courseId }) {
  const token = localStorage.getItem("auth_token");

  const [schedules, setSchedules] = useState([]);
  const [lessons, setLessons] = useState([]);
  const [teachers, setTeachers] = useState([]);

  const [selectedSchedule, setSelectedSchedule] = useState(null);
  const [existingLog, setExistingLog] = useState(null);
  const [showLogModal, setShowLogModal] = useState(false);

  const [form, setForm] = useState({
    lesson_id: "",
    teacher_id: "",
    planned_date: "",
    start_time: "",
    end_time: "",
  });

  const [logForm, setLogForm] = useState({
    conducted_date: "",
    completion_status: "COMPLETED",
    remarks: "",
    homework: "",
    next_topic: "",
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

  async function fetchModulesAndLessons() {
    try {
      const res = await http.get(`/courses/${courseId}/modules`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      const modules = res.data.modules || [];

      let allLessons = [];

      for (const module of modules) {
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
      alert("Failed to fetch lessons");
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

  async function openLog(schedule) {
    setSelectedSchedule(schedule);

    try {
      const res = await http.get(
        `/schedules/${schedule.id}/log`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      const log = res.data.log;

      setExistingLog(log);

      setLogForm({
        conducted_date: log.conducted_date.split("T")[0],
        completion_status: log.completion_status,
        remarks: log.remarks || "",
        homework: log.homework || "",
        next_topic: log.next_topic || "",
      });
    } catch (err) {
      setExistingLog(null);

      setLogForm({
        conducted_date: "",
        completion_status: "COMPLETED",
        remarks: "",
        homework: "",
        next_topic: "",
      });
    }

    setShowLogModal(true);
  }

  async function saveLog(e) {
    e.preventDefault();

    try {
      if (existingLog) {
        await http.put(
          `/course-logs/${existingLog.id}`,
          {
            completion_status:
              logForm.completion_status,
            remarks: logForm.remarks,
            homework: logForm.homework,
            next_topic: logForm.next_topic,
          },
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
      } else {
        await http.post(
          `/schedules/${selectedSchedule.id}/log`,
          logForm,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
      }

      alert("Log saved");

      setShowLogModal(false);
    } catch (err) {
      console.error(err);

      alert(
        err?.response?.data?.error ||
          "Failed to save log"
      );
    }
  }

  useEffect(() => {
    fetchSchedules();
    fetchModulesAndLessons();
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
              <th className="text-left p-2">Actions</th>
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

                <td className="p-2">
                  <button
                    onClick={() =>
                      openLog(schedule)
                    }
                    className="bg-indigo-600 text-white px-3 py-1 rounded"
                  >
                    Log
                  </button>
                </td>
              </tr>
            ))}

            {schedules.length === 0 && (
              <tr>
                <td
                  colSpan="7"
                  className="text-center p-4 text-gray-500"
                >
                  No schedules found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {showLogModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center">
          <div className="bg-white rounded shadow p-6 w-[700px]">
            <h2 className="text-xl font-semibold mb-4">
              {existingLog
                ? "Edit Session Log"
                : "Create Session Log"}
            </h2>

            <form
              onSubmit={saveLog}
              className="space-y-4"
            >
              {!existingLog && (
                <input
                  type="date"
                  value={logForm.conducted_date}
                  onChange={(e) =>
                    setLogForm({
                      ...logForm,
                      conducted_date:
                        e.target.value,
                    })
                  }
                  className="w-full border rounded p-2"
                  required
                />
              )}

              <select
                value={
                  logForm.completion_status
                }
                onChange={(e) =>
                  setLogForm({
                    ...logForm,
                    completion_status:
                      e.target.value,
                  })
                }
                className="w-full border rounded p-2"
              >
                <option value="COMPLETED">
                  COMPLETED
                </option>

                <option value="PARTIALLY_COMPLETED">
                  PARTIALLY_COMPLETED
                </option>

                <option value="CANCELLED">
                  CANCELLED
                </option>
              </select>

              <textarea
                placeholder="Remarks"
                value={logForm.remarks}
                onChange={(e) =>
                  setLogForm({
                    ...logForm,
                    remarks: e.target.value,
                  })
                }
                className="w-full border rounded p-2"
              />

              <textarea
                placeholder="Homework"
                value={logForm.homework}
                onChange={(e) =>
                  setLogForm({
                    ...logForm,
                    homework: e.target.value,
                  })
                }
                className="w-full border rounded p-2"
              />

              <textarea
                placeholder="Next Topic"
                value={logForm.next_topic}
                onChange={(e) =>
                  setLogForm({
                    ...logForm,
                    next_topic: e.target.value,
                  })
                }
                className="w-full border rounded p-2"
              />

              <div className="flex gap-2">
                <button
                  type="submit"
                  className="bg-green-600 text-white px-4 py-2 rounded"
                >
                  Save
                </button>

                <button
                  type="button"
                  onClick={() =>
                    setShowLogModal(false)
                  }
                  className="bg-gray-500 text-white px-4 py-2 rounded"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}