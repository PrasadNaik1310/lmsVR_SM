import { http } from "../services/http.js";

function withAuth(token) {
    console.log(`Auth token recived in withAuth() ${token}`)
  const config = token
    ? {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    : {};

  console.log("REQUEST CONFIG", config);

  return config;
}

export async function createCourse(payload, token) {
  const response = await http.post(
    "/courses",
    payload,
    withAuth(token)
  );
  
  return response.data;
}

export async function listCourses(params = {}, token) {
  const response = await http.get(
    "/courses",
    {
      ...withAuth(token),
      params,
    }
  );
  
  return response.data;
}

export async function getCourseDetails(courseId, token) {
  const response = await http.get(
    `/courses/${courseId}`,
    withAuth(token)
  );
  return response.data;
}

export async function updateCourse(courseId, payload, token) {
  const response = await http.put(
    `/courses/${courseId}`,
    payload,
    withAuth(token)
  );
  return response.data;
}

export async function publishCourse(courseId, token) {
  const response = await http.patch(
    `/courses/${courseId}/publish`,
    {},
    withAuth(token)
  );
  return response.data;
}

export async function generateInvite(courseId, token) {
  const response = await http.post(
    `/courses/${courseId}/invite`,
    {},
    withAuth(token)
  );
  return response.data;
}