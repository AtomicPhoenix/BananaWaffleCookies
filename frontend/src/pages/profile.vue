<template>
  <div class="page-container">
    <h2 class="page-title">Profile</h2>

    <!-- COMPLETION -->
    <div class="completion-section">
      <div class="completion-header">
        <span>Profile Completion</span>
        <span>{{ completionPercentage }}%</span>
      </div>
      <div class="progress-bar">
        <div
          class="progress-fill"
          :style="{ width: completionPercentage + '%' }"
        ></div>
      </div>
    </div>

    <div class="form-card">
      <!-- ================= BASIC INFO ================= -->
      <div class="section">
        <h3 class="section-title">Basic Information</h3>

        <div class="form-row">
          <input v-model="form.first_name" placeholder="First Name" />
          <input v-model="form.last_name" placeholder="Last Name" />
        </div>

        <input v-model="form.phone" placeholder="Phone" />

        <button @click="saveBasic">Save Basic Info</button>

        <p v-if="messages.basic.success" class="success">Saved!</p>
        <p v-if="messages.basic.error" class="error">
          {{ messages.basic.error }}
        </p>
      </div>

      <!-- ================= EDUCATION ================= -->
      <div class="section">
        <h3 class="section-title">Education</h3>

        <div class="form-row">
          <input v-model="newEducation.institution" placeholder="School" />
          <input v-model="newEducation.degree" placeholder="Degree" />
        </div>

        <div class="form-row">
          <input v-model="newEducation.field_of_study" placeholder="Field" />
          <input v-model="newEducation.start_date" type="date" />
          <input v-model="newEducation.end_date" type="date" />
        </div>

        <button @click="addEducation">Save Education</button>

        <p v-if="messages.education.success" class="success">Saved!</p>
        <p v-if="messages.education.error" class="error">
          {{ messages.education.error }}
        </p>

        <div
          v-for="edu in educationList"
          :key="edu.id"
          class="item-card"
        >
          <div v-if="editId !== edu.id" class="item-row">
            <div>
              <strong>{{ edu.institution }}</strong> — {{ edu.degree }}
              <p class="sub-text">{{ edu.field_of_study }}</p>
            </div>

            <div class="actions">
              <button @click="startEdit(edu)">Edit</button>
              <button @click="deleteEducation(edu.id)">Delete</button>
            </div>
          </div>

          <div v-else class="edit-row">
            <input v-model="editEducation.institution" />
            <input v-model="editEducation.degree" />
            <input v-model="editEducation.field_of_study" />
            <input v-model="editEducation.start_date" type="date" />
            <input v-model="editEducation.end_date" type="date" />

            <div class="actions">
              <button @click="updateEducation(edu.id)">Save</button>
              <button @click="cancelEdit">Cancel</button>
            </div>
          </div>
        </div>
      </div>

      <!-- ================= EMPLOYMENT ================= -->
      <div class="section">
        <h3 class="section-title">Employment</h3>

        <input v-model="newEmployment.organization" placeholder="Company" />
        <input v-model="newEmployment.title" placeholder="Role" />
        <input v-model="newEmployment.start_date" type="date" />
        <input v-model="newEmployment.end_date" type="date" />
        <input
          v-model="newEmployment.description"
          placeholder="Description"
        />

        <button @click="addEmployment">Add Employment</button>

        <p v-if="messages.employment.success" class="success">Saved!</p>
        <p v-if="messages.employment.error" class="error">
          {{ messages.employment.error }}
        </p>

        <div
          v-for="(emp, index) in employmentList"
          :key="emp.id"
          class="item-card"
        >
          <div v-if="editEmploymentId !== emp.id" class="item-row">
            <div>
              <strong>{{ emp.organization }}</strong> — {{ emp.title }}
              <p class="sub-text">{{ emp.description }}</p>
            </div>

            <div class="actions">
              <button @click="startEditEmployment(emp)">Edit</button>
              <button @click="moveEmploymentUp(index)">↑</button>
              <button @click="moveEmploymentDown(index)">↓</button>
              <button @click="deleteEmployment(emp.id)">Delete</button>
            </div>
          </div>

          <div v-else class="edit-row">
            <input v-model="editEmployment.organization" />
            <input v-model="editEmployment.title" />
            <input v-model="editEmployment.description" />

            <div class="actions">
              <button @click="updateEmployment(emp.id)">Save</button>
              <button @click="cancelEditEmployment">Cancel</button>
            </div>
          </div>
        </div>
      </div>

      <!-- ================= SKILLS ================= -->
      <div class="section">
        <h3 class="section-title">Skills</h3>

        <div class="form-row">
          <input v-model="newSkill.skill_name" placeholder="Skill" />
          <input
            v-model="newSkill.category"
            placeholder="Category (optional)"
          />

          <select v-model="newSkill.proficiency_label">
            <option value="">Proficiency</option>
            <option>Beginner</option>
            <option>Intermediate</option>
            <option>Advanced</option>
          </select>
        </div>

        <button @click="addSkill">Save Skill</button>

        <p v-if="messages.skills.success" class="success">Saved!</p>
        <p v-if="messages.skills.error" class="error">
          {{ messages.skills.error }}
        </p>

        <div
          v-for="(skill, index) in skillsList"
          :key="skill.id"
          class="item-card"
        >
          <div v-if="editSkillId !== skill.id" class="item-row">
            <div>
              <strong>{{ skill.skill_name }}</strong>
              <p class="sub-text">
                {{ skill.category }} • {{ skill.proficiency_label }}
              </p>
            </div>

            <div class="actions">
              <button @click="startEditSkill(skill)">Edit</button>
              <button @click="moveSkillUp(index)">↑</button>
              <button @click="moveSkillDown(index)">↓</button>
              <button @click="deleteSkill(skill.id)">Delete</button>
            </div>
          </div>

          <div v-else class="edit-row">
            <input v-model="editSkill.skill_name" />
            <input v-model="editSkill.category" />

            <select v-model="editSkill.proficiency_label">
              <option value="">Proficiency</option>
              <option>Beginner</option>
              <option>Intermediate</option>
              <option>Advanced</option>
            </select>

            <div class="actions">
              <button @click="updateSkill(skill.id)">Save</button>
              <button @click="cancelEditSkill">Cancel</button>
            </div>
          </div>
        </div>
      </div>

      <!-- ================= PROJECTS ================= -->
      <div class="section">
        <h3 class="section-title">Projects</h3>

        <input v-model="newProject.title" placeholder="Title" />
        <input v-model="newProject.description" placeholder="Description" />
        <input v-model="newProject.link" placeholder="Link" />

        <button @click="addProject">Add Project</button>

        <p v-if="messages.projects.success" class="success">Saved!</p>
        <p v-if="messages.projects.error" class="error">
          {{ messages.projects.error }}
        </p>

        <div
          v-for="(proj, index) in projectList"
          :key="proj.id"
          class="item-card"
        >
          <div v-if="editProjectId !== proj.id" class="item-row">
            <div>
              <strong>{{ proj.title }}</strong>
              <p class="sub-text">{{ proj.description }}</p>
            </div>

            <div class="actions">
              <button @click="startEditProject(proj)">Edit</button>
              <button @click="moveProjectUp(index)">↑</button>
              <button @click="moveProjectDown(index)">↓</button>
              <button @click="deleteProject(proj.id)">Delete</button>
            </div>
          </div>

          <div v-else class="edit-row">
            <input v-model="editProject.title" />
            <input v-model="editProject.description" />

            <div class="actions">
              <button @click="updateProject(proj.id)">Save</button>
              <button @click="cancelEditProject">Cancel</button>
            </div>
          </div>
        </div>
      </div>

      <!-- ================= PREFERENCES ================= -->
      <div class="section">
        <h3 class="section-title">Preferences</h3>

        <input
          v-model="preferences.preferred_role"
          placeholder="Target Roles"
        />
        <input
          v-model="preferences.location"
          placeholder="Preferred Location"
        />

        <select v-model="preferences.work_mode">
          <option value="">Work Mode</option>
          <option>Remote</option>
          <option>Hybrid</option>
          <option>On-site</option>
        </select>

        <input
          v-model="preferences.salary"
          type="number"
          placeholder="Desired Salary"
        />

        <button @click="savePreferences">Save Preferences</button>

        <p v-if="messages.preferences.success" class="success">
          Saved!
        </p>
        <p v-if="messages.preferences.error" class="error">
          {{ messages.preferences.error }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed, onMounted } from 'vue'

const form = reactive({
  first_name: '',
  last_name: '',
  phone: ''
})

const preferences = reactive({
  preferred_role: '',
  location: '',
  work_mode: '',
  salary: ''
})

const educationList = ref([])
const skillsList = ref([])
const employmentList = ref([])
const projectList = ref([])

const newEducation = reactive({
  institution: '',
  degree: '',
  field_of_study: '',
  start_date: '',
  end_date: ''
})

const newSkill = reactive({
  skill_name: '',
  category: '',
  proficiency_label: ''
})

const newEmployment = reactive({
  organization: '',
  title: '',
  start_date: '',
  end_date: '',
  description: ''
})

const newProject = reactive({
  title: '',
  description: '',
  link: ''
})

const editId = ref(null)
const editEducation = reactive({
  institution: '',
  degree: '',
  field_of_study: '',
  start_date: '',
  end_date: ''
})

const editSkillId = ref(null)
const editSkill = reactive({
  skill_name: '',
  category: '',
  proficiency_label: ''
})

const editEmploymentId = ref(null)
const editEmployment = reactive({
  organization: '',
  title: '',
  start_date: '',
  end_date: '',
  description: ''
})

const editProjectId = ref(null)
const editProject = reactive({
  title: '',
  description: '',
  link: ''
})

const messages = reactive({
  basic: { success: false, error: '' },
  education: { success: false, error: '' },
  skills: { success: false, error: '' },
  employment: { success: false, error: '' },
  projects: { success: false, error: '' },
  preferences: { success: false, error: '' }
})

const completionPercentage = computed(() => {
  const fields = [
    form.first_name,
    form.last_name,
    form.phone,
    preferences.preferred_role,
    preferences.location,
    preferences.work_mode,
    preferences.salary
  ]

  const filled = fields.filter(v => v && v.toString().trim() !== '').length

  return Math.round((filled / fields.length) * 100)
})


onMounted(() => {
  getProfile()
  getEducation()
  getSkills()
  getEmployment()
  getProjects()
})

async function getProfile() {
  try {
    const res = await fetch('/api/profile')
    if (res.ok) {
      const data = await res.json()
      Object.assign(form, data)
      Object.assign(preferences, data.preferences || {})
    }
  } catch (err) {
    console.error(err)
  }
}


async function getEducation() {
  const res = await fetch('/api/profile/education')
  if (res.ok) educationList.value = await res.json()
}

async function getSkills() {
  const res = await fetch('/api/profile/skills')
  if (res.ok) skillsList.value = await res.json()
}

async function saveSkillOrder() {
  try {
    await fetch('/api/profile/skills/reorder', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(
        skillsList.value.map((skill, index) => ({
          id: skill.id,
          sort_order: index
        }))
      )
    })
  } catch (err) {
    console.error('Failed to save skill order', err)
  }
}

async function getEmployment() {
  const res = await fetch('/api/profile/employment')
  if (res.ok) employmentList.value = await res.json()
}

async function saveEmploymentOrder() {
  try {
    await fetch('/api/profile/employment/reorder', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(
        employmentList.value.map((emp, index) => ({
          id: emp.id,
          sort_order: index
        }))
      )
    })
  } catch (err) {
    console.error('Failed to save employment order', err)
  }
}

async function getProjects() {
  const res = await fetch('/api/profile/projects')
  if (res.ok) projectList.value = await res.json()
}

async function saveProjectOrder() {
  try {
    await fetch('/api/profile/projects/reorder', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(
        projectList.value.map((proj, index) => ({
          id: proj.id,
          sort_order: index
        }))
      )
    })
  } catch (err) {
    console.error('Failed to save project order', err)
  }
}

function reset(section) {
  messages[section].success = false
  messages[section].error = ''
}

async function saveBasic() {
  reset('basic')

  if (!form.first_name || !form.last_name) {
    messages.basic.error = 'First and last name required'
    return
  }

  try {
    const res = await fetch('/api/profile/basic', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })

    if (res.ok) messages.basic.success = true
    else messages.basic.error = 'Save failed'
  } catch {
    messages.basic.error = 'Server error'
  }
}

async function savePreferences() {
  reset('preferences')

  if (!preferences.preferred_role || !preferences.location) {
    messages.preferences.error = 'Roles and location required'
    return
  }

  try {
    const res = await fetch('/api/profile/preferences', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(preferences)
    })

    if (res.ok) messages.preferences.success = true
    else messages.preferences.error = 'Save failed'
  } catch {
    messages.preferences.error = 'Server error'
  }
}

async function addEducation() {
  reset('education')

  if (
    !newEducation.institution ||
    !newEducation.degree ||
    !newEducation.field_of_study ||
    !newEducation.start_date ||
    !newEducation.end_date
  ) {
    messages.education.error = 'All fields required'
    return
  }

  if (newEducation.end_date < newEducation.start_date) {
    messages.education.error = 'End date must be after start date'
    return
  }

  try {
    const res = await fetch('/api/profile/education', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newEducation)
    })

    if (res.ok) {
      messages.education.success = true
      Object.keys(newEducation).forEach(k => (newEducation[k] = ''))
      getEducation()
    } else {
      messages.education.error = 'Save failed'
    }
  } catch {
    messages.education.error = 'Server error'
  }
}

async function deleteEducation(id) {
  await fetch(`/api/profile/education/${id}`, { method: 'DELETE' })
  getEducation()
}

function startEdit(edu) {
  editId.value = edu.id
  Object.assign(editEducation, edu)
}

function cancelEdit() {
  editId.value = null
}

async function updateEducation(id) {
  reset('education')

  if (
    !editEducation.institution ||
    !editEducation.degree ||
    !editEducation.field_of_study ||
    !editEducation.start_date ||
    !editEducation.end_date
  ) {
    messages.education.error = 'All fields required'
    return
  }

  if (editEducation.end_date < editEducation.start_date) {
    messages.education.error = 'End date must be after start date'
    return
  }

  try {
    const res = await fetch(`/api/profile/education/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(editEducation)
    })

    if (res.ok) {
      messages.education.success = true
      editId.value = null
      getEducation()
    } else {
      messages.education.error = 'Update failed'
    }
  } catch {
    messages.education.error = 'Server error'
  }
}

// ================= SKILLS =================

function startEditSkill(skill) {
  editSkillId.value = skill.id
  Object.assign(editSkill, skill)
}

function cancelEditSkill() {
  editSkillId.value = null
}

async function updateSkill(id) {
  reset('skills')

  if (!editSkill.skill_name) {
    messages.skills.error = 'Skill name required'
    return
  }

  try {
    const res = await fetch(`/api/profile/skills/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(editSkill)
    })

    if (res.ok) {
      messages.skills.success = true
      editSkillId.value = null
      getSkills()
    } else {
      messages.skills.error = 'Update failed'
    }
  } catch {
    messages.skills.error = 'Server error'
  }
}

async function addSkill() {
  reset('skills')

  if (!newSkill.skill_name) {
    messages.skills.error = 'Skill name required'
    return
  }

  try {
    const res = await fetch('/api/profile/skills', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newSkill)
    })

    if (res.ok) {
      messages.skills.success = true
      Object.keys(newSkill).forEach(k => (newSkill[k] = ''))
      getSkills()
    } else {
      messages.skills.error = 'Save failed'
    }
  } catch {
    messages.skills.error = 'Server error'
  }
}

async function deleteSkill(id) {
  await fetch(`/api/profile/skills/${id}`, { method: 'DELETE' })
  getSkills()
}

async function moveSkillUp(index) {
  if (index === 0) return
  const temp = skillsList.value[index]
  skillsList.value[index] = skillsList.value[index - 1]
  skillsList.value[index - 1] = temp
  await saveSkillOrder()
}

async function moveSkillDown(index) {
  if (index === skillsList.value.length - 1) return
  const temp = skillsList.value[index]
  skillsList.value[index] = skillsList.value[index + 1]
  skillsList.value[index + 1] = temp
  await saveSkillOrder()
}

async function addEmployment() {
  reset('employment')

  if (!newEmployment.organization || !newEmployment.title) {
    messages.employment.error = 'Company and role required'
    return
  }

  const res = await fetch('/api/profile/employment', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(newEmployment)
  })

  if (res.ok) {
    messages.employment.success = true
    Object.keys(newEmployment).forEach(k => (newEmployment[k] = ''))
    getEmployment()
  } else {
    messages.employment.error = 'Save failed'
  }
}

function startEditEmployment(emp) {
  editEmploymentId.value = emp.id
  Object.assign(editEmployment, emp)
}

function cancelEditEmployment() {
  editEmploymentId.value = null
}

async function updateEmployment(id) {
  reset('employment')

  const res = await fetch(`/api/profile/employment/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(editEmployment)
  })

  if (res.ok) {
    messages.employment.success = true
    editEmploymentId.value = null
    getEmployment()
  } else {
    messages.employment.error = 'Update failed'
  }
}

async function deleteEmployment(id) {
  await fetch(`/api/profile/employment/${id}`, { method: 'DELETE' })
  getEmployment()
}

async function addProject() {
  reset('projects')

  if (!newProject.title) {
    messages.projects.error = 'Title required'
    return
  }

  const res = await fetch('/api/profile/projects', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(newProject)
  })

  if (res.ok) {
    messages.projects.success = true
    Object.keys(newProject).forEach(k => (newProject[k] = ''))
    getProjects()
  } else {
    messages.projects.error = 'Save failed'
  }
}

function startEditProject(p) {
  editProjectId.value = p.id
  Object.assign(editProject, p)
}

function cancelEditProject() {
  editProjectId.value = null
}

async function updateProject(id) {
  reset('projects')

  const res = await fetch(`/api/profile/projects/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(editProject)
  })

  if (res.ok) {
    messages.projects.success = true
    editProjectId.value = null
    getProjects()
  } else {
    messages.projects.error = 'Update failed'
  }
}

async function deleteProject(id) {
  await fetch(`/api/profile/projects/${id}`, { method: 'DELETE' })
  getProjects()
}

async function moveEmploymentUp(index) {
  if (index === 0) return
  const temp = employmentList.value[index]
  employmentList.value[index] = employmentList.value[index - 1]
  employmentList.value[index - 1] = temp
  await saveEmploymentOrder()
}

async function moveEmploymentDown(index) {
  if (index === employmentList.value.length - 1) return
  const temp = employmentList.value[index]
  employmentList.value[index] = employmentList.value[index + 1]
  employmentList.value[index + 1] = temp
  await saveEmploymentOrder()
}

async function moveProjectUp(index) {
  if (index === 0) return
  const temp = projectList.value[index]
  projectList.value[index] = projectList.value[index - 1]
  projectList.value[index - 1] = temp
  await saveProjectOrder()
}

async function moveProjectDown(index) {
  if (index === projectList.value.length - 1) return
  const temp = projectList.value[index]
  projectList.value[index] = projectList.value[index + 1]
  projectList.value[index + 1] = temp
  await saveProjectOrder()
}
</script>

<style scoped src="@/assets/css/profile.css"></style>