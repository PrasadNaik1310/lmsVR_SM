import React from "react";

export default function CourseSchedule({ courseId }) {
  return (
    <div className="bg-white shadow rounded p-6">
      <h2 className="text-xl font-semibold mb-2">Schedule</h2>
      <p className="text-sm text-gray-600 mb-4">Course ID: {courseId}</p>

      <div className="border rounded p-4 text-gray-600">
        Course schedule placeholder — implementation coming soon.
      </div>
    </div>
  );
}
