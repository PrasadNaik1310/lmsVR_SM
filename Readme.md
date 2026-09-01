# LMS Management System
 ## For technical implementation , PLease refer LMS_AWS_Deployment_Manual.md file 
## Overview

This project is a centralized Learning Management System (LMS) designed to manage academic, administrative, and student-related workflows within a single platform.

The V1 implementation establishes the core administrative, academic, authentication, authorization, and backend foundation required for the LMS. The system currently includes admissions, course management, academic sessions, batches, internal user management, scheduling, and role-based access control.

---

# Project Objectives

The platform is being developed to:

* Simplify student admissions and onboarding
* Centralize course and academic management
* Improve communication between teachers, students, and administrators
* Manage live classes and attendance digitally
* Track enrollments and payment status
* Reduce manual administrative work
* Provide a centralized platform for LMS administration

---
# Current Status

## V1 Foundation — Completed

The V1 foundation of the LMS has been implemented and deployed to a production environment.

### Completed

* System architecture and database design
* PostgreSQL database integration
* GORM-based data models
* JWT authentication
* bcrypt password hashing
* Role-based access control (RBAC)
* Permission-based API authorization
* Admissions workflow
* Course management
* Course modules and lessons
* Course scheduling and logging
* Academic session management
* Batch management
* Internal team/user management
* Student, teacher and admin account creation
* Frontend administration interfaces
* REST API integration
* Database migrations and seed data
* Docker-based development/deployment setup
* Production deployment

### V1 Scope

The current version establishes the core administrative and backend foundation of the LMS. Student-facing workflows and additional product modules are planned for subsequent development phases.
# Core Modules

## 1. Authentication & Authorization (Completed)

Provides authentication and role-based authorization across the platform.

### Features

* JWT-based authentication
* Password hashing using bcrypt
* Role-based access control (RBAC)
* Permission-based API authorization
* Role-permission mappings
* Protected API routes

---

## 2. Admissions Module (Completed)

Handles the admission workflow from enquiry through student account creation.

### Features

* Student enquiries
* Application management
* Application approval/rejection
* Student account creation
* Temporary password generation
* Student onboarding

---

## 3. Manage Company Module (Completed)

Handles organizational and academic management for the LMS.

### Features

* Academic session management
* Course/package management
* Batch management
* Internal team management
* Student, teacher and admin account creation
* Course-to-academic-session assignment
* Seat tracking
* Batch and session listing

---

## 4. Course Management Module (Completed)

Provides course creation and management functionality.

### Features

* Course creation
* Course listing
* Course details
* Course updates
* Course publishing
* Course invite link generation
* Course/module relationships

---

## 5. Course Modules & Lessons (Completed)

Provides hierarchical course content management.

### Features

* Module creation
* Module listing
* Module updates
* Module deletion
* Lesson creation
* Lesson listing
* Lesson details
* Lesson updates
* Lesson deletion

---

## 6. Course Planning & Logging (Completed)

Provides scheduling and academic session logging functionality.

### Features

* Course schedule creation
* Schedule listing
* Schedule details
* Schedule updates
* Schedule deletion
* Session logging
* Course log updates
* Recurring schedules

---

## 7. Internal Team & User Management (Completed)

Provides administrative user management for the organization.

### Features

* Create student accounts
* Create teacher accounts
* Create admin accounts
* List company users
* Role assignment
* Secure password hashing
* Role-based access control

---

## 8. Academic Session Management (Completed)

Provides management of academic sessions and their associated courses.

### Features

* Create academic sessions
* List academic sessions
* View academic session details
* Update academic sessions
* Delete academic sessions
* Assign courses to academic sessions
* List courses associated with an academic session

---

## 9. Batch Management (Completed)

Provides batch management under courses.

### Features

* Create batches
* List batches by course
* View batch details
* Batch seat management

---

# Modules Planned for Future Development

The following modules are part of the broader LMS product scope but are not part of the current V1 implementation.

## Membership Module

Planned functionality:

* Payment tracking
* Grace period handling
* Membership status management
* Enrollment tracking

---

## CRM Dashboard

Planned functionality:

* Total students
* Active courses
* Outstanding fees
* Live classes
* Recent activities
* Notifications
* Quick management actions

---

## Student Course Dashboard

Planned functionality:

* Student-facing course dashboard
* Access enrolled courses
* Access lessons
* View upcoming sessions
* Student learning experience

---

## Live Sessions Module

Planned functionality:

* Online class management
* Attendance tracking
* Notifications
* Session management

---

## Future Learning Features

Planned enhancements:

* Quizzes
* Assignments
* Reports & analytics
* Certificates
* Discussion forums

---

# User Roles

| Role | Responsibilities |
|---|---|
| Super Admin | Manage the complete platform |
| Admin | Manage students, courses, and announcements |
| Teacher | Create courses, lessons, and manage attendance |
| Finance Manager | Manage payments and memberships |
| Student | Access courses and attend sessions |

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