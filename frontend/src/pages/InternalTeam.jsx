import { useState } from "react";
import { createCompanyUser } from "../api/manageCompany.js";
import {
  listCompanyUsers,
} from "../api/manageCompany.js";
import { useEffect } from "react";
export default function InternalTeam() {
  const token = localStorage.getItem("auth_token");

  const [showModal, setShowModal] = useState(false);
  const [loading, setLoading] = useState(false);
    const [members, setMembers] = useState([]);
const [loadingMembers, setLoadingMembers] = useState(false);


  const [form, setForm] = useState({
    first_name: "",
    last_name: "",
    email: "",
    phone: "",
    password: "",
    role: "teacher",

    // Student
    enrollment_number: "",
    date_of_birth: "",
    address: "",
    admission_date: "",

    // Teacher
    specialization: "",
    bio: "",
    joining_date: "",
  });

  function handleChange(e) {
    setForm((prev) => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  }

  function handleRoleChange(e) {
    setForm((prev) => ({
      ...prev,
      role: e.target.value,
    }));
  }

  function resetForm() {
    setForm({
      first_name: "",
      last_name: "",
      email: "",
      phone: "",
      password: "",
      role: "teacher",

      enrollment_number: "",
      date_of_birth: "",
      address: "",
      admission_date: "",

      specialization: "",
      bio: "",
      joining_date: "",
    });
  }
useEffect(() => {
  fetchMembers();
}, []);
  async function handleCreateMember(e) {
    e.preventDefault();

    try {
      setLoading(true);

      const payload = {
        first_name: form.first_name,
        last_name: form.last_name,
        email: form.email,
        phone: form.phone,
        password: form.password,
        role: form.role,
      };

      if (form.role === "student") {
        payload.enrollment_number = form.enrollment_number;
        payload.date_of_birth = form.date_of_birth;
        payload.address = form.address;
        payload.admission_date = form.admission_date;
      }

      if (form.role === "teacher") {
        payload.specialization = form.specialization;
        payload.bio = form.bio;
        payload.joining_date = form.joining_date;
      }

      await createCompanyUser(payload, token);

      alert("Member created successfully");

      resetForm();
      setShowModal(false);
      await fetchMembers();
    } catch (err) {
      console.error(err);

      alert(
        err?.response?.data?.error ||
          "Failed to create member"
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="p-6 space-y-8">

      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">
            Internal Team
          </h1>

          <p className="text-slate-500 mt-1">
            Manage students, teachers and administrators.
          </p>
        </div>

        <button
          onClick={() => setShowModal(true)}
          className="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700"
        >
          + Add Member
        </button>
      </div>

      <div className="bg-white rounded-lg shadow p-6">

  <div className="flex justify-between items-center mb-4">
    <h2 className="text-xl font-semibold">
      Team Members
    </h2>

    <button
      onClick={fetchMembers}
      disabled={loadingMembers}
      className="border px-4 py-2 rounded"
    >
      {loadingMembers ? "Refreshing..." : "Refresh"}
    </button>
  </div>

  {loadingMembers ? (
    <p>Loading team members...</p>
  ) : members.length === 0 ? (
    <p className="text-slate-500">
      No team members found.
    </p>
  ) : (
    <table className="w-full border-collapse">
      <thead>
        <tr className="border-b">
          <th className="text-left p-3">Name</th>
          <th className="text-left p-3">Email</th>
          <th className="text-left p-3">Phone</th>
          <th className="text-left p-3">Role</th>
          <th className="text-left p-3">Status</th>
        </tr>
      </thead>

      <tbody>
        {members.map((member) => (
          <tr
            key={member.id}
            className="border-b"
          >
            <td className="p-3">
              {member.first_name} {member.last_name}
            </td>

            <td className="p-3">
              {member.email}
            </td>

            <td className="p-3">
              {member.phone || "-"}
            </td>

            <td className="p-3 capitalize">
              {member.role}
            </td>

            <td className="p-3">
              {member.is_active ? "Active" : "Inactive"}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )}
</div>

      {/* Add Member Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">

          <div className="bg-white rounded-xl shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto p-6">

            <div className="flex justify-between items-center mb-6">
              <h2 className="text-2xl font-semibold">
                Add Member
              </h2>

              <button
                type="button"
                onClick={() => {
                  resetForm();
                  setShowModal(false);
                }}
                className="text-slate-500 hover:text-slate-800 text-xl"
              >
                ×
              </button>
            </div>

            <form
              onSubmit={handleCreateMember}
              className="space-y-6"
            >

              {/* Common fields */}
              <div>
                <h3 className="font-semibold mb-3">
                  Account Details
                </h3>

                <div className="grid grid-cols-2 gap-4">

                  <input
                    name="first_name"
                    placeholder="First Name"
                    value={form.first_name}
                    onChange={handleChange}
                    required
                    className="border p-2 rounded"
                  />

                  <input
                    name="last_name"
                    placeholder="Last Name"
                    value={form.last_name}
                    onChange={handleChange}
                    className="border p-2 rounded"
                  />

                  <input
                    name="email"
                    type="email"
                    placeholder="Email"
                    value={form.email}
                    onChange={handleChange}
                    required
                    className="border p-2 rounded"
                  />

                  <input
                    name="phone"
                    placeholder="Phone"
                    value={form.phone}
                    onChange={handleChange}
                    required
                    className="border p-2 rounded"
                  />

                  <input
                    name="password"
                    type="password"
                    placeholder="Password"
                    value={form.password}
                    onChange={handleChange}
                    required
                    minLength={6}
                    className="border p-2 rounded"
                  />

                  <select
                    name="role"
                    value={form.role}
                    onChange={handleRoleChange}
                    className="border p-2 rounded"
                  >
                    <option value="teacher">
                      Teacher
                    </option>

                    <option value="student">
                      Student
                    </option>

                    <option value="admin">
                      Admin
                    </option>
                  </select>

                </div>
              </div>

              {/* Student fields */}
              {form.role === "student" && (
                <div>
                  <h3 className="font-semibold mb-3">
                    Student Details
                  </h3>

                  <div className="grid grid-cols-2 gap-4">

                    <input
                      name="enrollment_number"
                      placeholder="Enrollment Number"
                      value={form.enrollment_number}
                      onChange={handleChange}
                      className="border p-2 rounded"
                    />

                    <input
                      name="date_of_birth"
                      type="date"
                      value={form.date_of_birth}
                      onChange={handleChange}
                      className="border p-2 rounded"
                    />

                    <input
                      name="admission_date"
                      type="date"
                      value={form.admission_date}
                      onChange={handleChange}
                      required
                      className="border p-2 rounded"
                    />

                    <textarea
                      name="address"
                      placeholder="Address"
                      value={form.address}
                      onChange={handleChange}
                      className="border p-2 rounded col-span-2"
                      rows={3}
                    />

                  </div>
                </div>
              )}

              {/* Teacher fields */}
              {form.role === "teacher" && (
                <div>
                  <h3 className="font-semibold mb-3">
                    Teacher Details
                  </h3>

                  <div className="grid grid-cols-2 gap-4">

                    <input
                      name="specialization"
                      placeholder="Specialization"
                      value={form.specialization}
                      onChange={handleChange}
                      className="border p-2 rounded"
                    />

                    <input
                      name="joining_date"
                      type="date"
                      value={form.joining_date}
                      onChange={handleChange}
                      className="border p-2 rounded"
                    />

                    <textarea
                      name="bio"
                      placeholder="Bio"
                      value={form.bio}
                      onChange={handleChange}
                      className="border p-2 rounded col-span-2"
                      rows={4}
                    />

                  </div>
                </div>
              )}

              {/* Buttons */}
              <div className="flex justify-end gap-3 pt-4 border-t">

                <button
                  type="button"
                  onClick={() => {
                    resetForm();
                    setShowModal(false);
                  }}
                  className="border px-4 py-2 rounded"
                  disabled={loading}
                >
                  Cancel
                </button>

                <button
                  type="submit"
                  disabled={loading}
                  className="bg-indigo-600 text-white px-4 py-2 rounded disabled:opacity-50"
                >
                  {loading
                    ? "Creating..."
                    : "Create Member"}
                </button>

              </div>

            </form>
          </div>
        </div>
      )}
    </div>
  );
  async function fetchMembers() {
  try {
    setLoadingMembers(true);

    const data = await listCompanyUsers(token);

    console.log("COMPANY USERS:", data);

    setMembers(data.users || []);
  } catch (err) {
    console.error("Failed to load company users:", err);

    alert(
      err?.response?.data?.error ||
        "Failed to load team members"
    );
  } finally {
    setLoadingMembers(false);
  }
}
}
