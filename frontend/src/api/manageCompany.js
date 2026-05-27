import { http } from "../services/http.js";

function withAuth(token) {
  return token
    ? {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    : {};
}

export async function listBatchesByCourse(courseId, params = {}, token) {
  const response = await http.get(`/company/courses/${courseId}/batches`, {
    ...withAuth(token),
    params,
  });
  return response.data;
}

export async function createBatchForCourse(courseId, payload, token) {
  const response = await http.post(`/company/courses/${courseId}/batches`, payload, withAuth(token));
  return response.data;
}

export async function getBatchDetails(courseId, batchId, token) {
  const response = await http.get(`/company/courses/${courseId}/batches/${batchId}`, withAuth(token));
  return response.data;
}

export async function listCoursesBySession(sessionId, params = {}, token) {
  const response = await http.get(`/company/sessions/${sessionId}/courses`, {
    ...withAuth(token),
    params,
  });
  return response.data;
}

export async function assignCourseToSession(sessionId, courseId, token) {
  const response = await http.put(
    `/company/sessions/${sessionId}/courses/${courseId}/assign`,
    undefined,
    withAuth(token),
  );
  return response.data;
}