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

// Enquiry endpoints
export async function createEnquiry(payload, token) {
  const response = await http.post(`/admissions/enquiries`, payload, withAuth(token));
  return response.data;
}

export async function listEnquiries(params = {}, token) {
  const response = await http.get(`/admissions/enquiries`, {
    ...withAuth(token),
    params,
  });
  return response.data;
}

export async function getEnquiry(enquiryId, token) {
  const response = await http.get(`/admissions/enquiries/${enquiryId}`, withAuth(token));
  return response.data;
}

export async function updateEnquiryStatus(enquiryId, payload, token) {
  const response = await http.patch(
    `/admissions/enquiries/${enquiryId}/status`,
    payload,
    withAuth(token)
  );
  return response.data;
}

// Application endpoints
export async function createApplication(payload, token) {
  const response = await http.post(`/admissions/applications`, payload, withAuth(token));
  return response.data;
}

export async function listApplications(params = {}, token) {
  const response = await http.get(`/admissions/applications`, {
    ...withAuth(token),
    params,
  });
  return response.data;
}

export async function getApplication(applicationId, token) {
  const response = await http.get(`/admissions/applications/${applicationId}`, withAuth(token));
  return response.data;
}

export async function approveApplication(applicationId, payload, token) {
  const response = await http.patch(
    `/admissions/applications/${applicationId}/approve`,
    payload,
    withAuth(token)
  );
  return response.data;
}

export async function rejectApplication(applicationId, payload, token) {
  const response = await http.patch(
    `/admissions/applications/${applicationId}/reject`,
    payload,
    withAuth(token)
  );
  return response.data;
}
