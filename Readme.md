# LMS Management System

## Overview

This project is a centralized Learning Management System (LMS) designed to manage academic, administrative, and student-related workflows within a single platform.

The system combines admissions management, course management, live classes, enrollment tracking, attendance, announcements, and payment management into one unified solution.

---

# Project Objectives

The platform is being developed to:

* Simplify student admissions and onboarding
* Centralize course and academic management
* Improve communication between teachers, students, and administrators
* Manage live classes and attendance digitally
* Track enrollments and payment status
* Reduce manual administrative work

---

# Core Modules

## 1. Admissions Module

Handles the complete admission workflow from enquiry to student onboarding.

### Features

* Student enquiries
* Application management
* Admission approval/rejection
* Student account creation

---

## 2. Manage Company Module

Handles organizational and academic management.

### Features

* Academic sessions
* Batch management
* Internal team management
* Course/package management
* Seat tracking

---

## 3. Membership Module

Manages student payments and access control.

### Features

* Payment tracking
* Grace period handling
* Membership status management
* Enrollment tracking

---

## 4. CRM Dashboard

Central dashboard for administrators.

### Dashboard Includes

* Total students
* Active courses
* Outstanding fees
* Live classes
* Recent activities
* Notifications
* Quick management actions

---

## 5. Course Creation Module

Allows teachers and administrators to create and manage courses.

### Features

* Course creation
* Module and lesson management
* PDF/resource uploads
* Course publishing
* Invite link generation

---

## 6. Course Planning & Logging

Used for scheduling and tracking academic sessions.

### Features

* Session planning
* Timeline scheduling
* Recurring schedules
* Session logs

---

## 7. Course List Dashboard

Displays all enrolled courses for students.

### Features

* Access enrolled courses
* Access lessons
* View upcoming sessions

---

## 8. Live Sessions Module

Handles online class sessions.

### Features

* Session scheduling
* Attendance tracking
* Notifications
* Session management

---

# User Roles

| Role            | Responsibilities                               |
| --------------- | ---------------------------------------------- |
| Super Admin     | Manage the complete platform                   |
| Admin           | Manage students, courses, and announcements    |
| Teacher         | Create courses, lessons, and manage attendance |
| Finance Manager | Manage payments and memberships                |
| Student         | Access courses and attend sessions             |

---

# High-Level Workflow

```text
Enquiry
    ↓
Application
    ↓
Admission Approval
    ↓
Student Account Creation
    ↓
Course Enrollment
    ↓
Membership & Payment
    ↓
Access Lessons & Live Sessions
    ↓
Attendance & Progress Tracking
```

---

# Key Features

* Centralized LMS management
* Admission and enrollment workflow
* Live session management
* Attendance tracking
* Payment and membership management
* Course planning and scheduling
* Announcement and notification system
* Role-based access management

---

# Future Scope

Planned future enhancements:

* Quizzes
* Assignments
* Reports & analytics
* Certificates
* Discussion forums

---

# Current Status

```text
Architecture & Database Design Completed
Gorm modules setup , basic auth setup (JWT)
```

---

# Documentation Included

* ER Diagram
* Database Schema
* High-Level Workflow
* Module-wise Workflows
* Role Definitions

---

# Version

Version: v1.0
Last Updated: May 2026
