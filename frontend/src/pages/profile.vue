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
        <div class="progress-fill" :style="{ width: completionPercentage + '%' }"></div>
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
        <p v-if="messages.basic.error" class="error">{{ messages.basic.error }}</p>
      </div>

      <!-- ================= EDUCATION ================= -->
      <div class="section">
        <h3 class="section-title">Education</h3>

        <div class="form-row">
          <input v-model="newEducation.school" placeholder="School" />
          <input v-model="newEducation.degree" placeholder="Degree" />
        </div>

        <div class="form-row">
          <input v-model="newEducation.field" placeholder="Field" />
          <input v-model="newEducation.start_date" type="date" />
          <input v-model="newEducation.end_date" type="date" />
        </div>

        <button @click="addEducation">Save Education</button>
        <p v-if="messages.education.success" class="success">Saved!</p>
        <p v-if="messages.education.error" class="error">{{ messages.education.error }}</p>

        <div v-for="edu in educationList" :key="edu.id" class="item-card">
          <div v-if="editId !== edu.id" class="item-row">
            <div>
              <strong>{{ edu.school }}</strong> — {{ edu.degree }}
              <p class="sub-text">{{ edu.field }}</p>
            </div>
            <div class="actions">
              <button @click="startEdit(edu)">Edit</button>
              <button @click="deleteEducation(edu.id)">Delete</button>
            </div>
          </div>

          <div v-else class="edit-row">
            <input v-model="editEducation.school" />
            <input v-model="editEducation.degree" />
            <input v-model="editEducation.field" />
            <input v-model="editEducation.start_date" type="date" />
            <input v-model="editEducation.end_date" type="date" />

            <div class="actions">
              <button @click="updateEducation(edu.id)">Save</button>
              <button @click="cancelEdit">Cancel</button>
            </div>
          </div>
        </div>
      </div>

      <!-- ================= SKILLS ================= -->
      <div class="section">
        <h3 class="section-title">Skills</h3>

        <div class="form-row">
          <input v-model="newSkill.name" placeholder="Skill" />
          <input v-model="newSkill.category" placeholder="Category (optional)" />
          <select v-model="newSkill.proficiency">
            <option value="">Proficiency</option>
            <option>Beginner</option>
            <option>Intermediate</option>
            <option>Advanced</option>
          </select>
        </div>

        <button @click="addSkill">Save Skill</button>
        <p v-if="messages.skills.success" class="success">Saved!</p>
        <p v-if="messages.skills.error" class="error">{{ messages.skills.error }}</p>

        <div v-for="(skill, index) in skillsList" :key="skill.id" class="item-card">
          <div v-if="editSkillId !== skill.id" class="item-row">
            <div>
              <strong>{{ skill.name }}</strong>
              <p class="sub-text">{{ skill.category }} • {{ skill.proficiency }}</p>
            </div>

            <div class="actions">
              <button @click="startEditSkill(skill)">Edit</button>
              <button @click="moveSkillUp(index)">↑</button>
              <button @click="moveSkillDown(index)">↓</button>
              <button @click="deleteSkill(skill.id)">Delete</button>
            </div>
          </div>

          <div v-else class="edit-row">
            <input v-model="editSkill.name" />
            <input v-model="editSkill.category" />
            <select v-model="editSkill.proficiency">
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

      <!-- ================= PREFERENCES ================= -->
      <div class="section">
        <h3 class="section-title">Preferences</h3>

        <input v-model="preferences.target_roles" placeholder="Target Roles" />
        <input v-model="preferences.location" placeholder="Preferred Location" />

        <select v-model="preferences.work_mode">
          <option value="">Work Mode</option>
          <option>Remote</option>
          <option>Hybrid</option>
          <option>On-site</option>
        </select>

        <input v-model="preferences.salary" type="number" placeholder="Desired Salary" />

        <button @click="savePreferences">Save Preferences</button>
        <p v-if="messages.preferences.success" class="success">Saved!</p>
        <p v-if="messages.preferences.error" class="error">{{ messages.preferences.error }}</p>
      </div>

    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed, onMounted } from 'vue'

const form = reactive({ first_name: '', last_name: '', phone: '' })
const preferences = reactive({ target_roles: '', location: '', work_mode: '', salary: '' })

const educationList = ref([])
const skillsList = ref([])

const newEducation = reactive({ school: '', degree: '', field: '', start_date: '', end_date: '' })
const newSkill = reactive({ name: '', category: '', proficiency: '' })

const editId = ref(null)
const editEducation = reactive({ school: '', degree: '', field: '', start_date: '', end_date: '' })

const editSkillId = ref(null)
const editSkill = reactive({ name: '', category: '', proficiency: '' })

const messages = reactive({
  basic: { success: false, error: '' },
  education: { success: false, error: '' },
  skills: { success: false, error: '' },
  preferences: { success: false, error: '' }
})

const completionPercentage = computed(() => {
  const total = [...Object.values(form), preferences.target_roles, preferences.location]
  const filled = total.filter(v => v && v.toString().trim() !== '').length
  return Math.round((filled / total.length) * 100)
})

onMounted(() => {
  getProfile()
  getEducation()
  getSkills()
})

async function getProfile() {
  const res = await fetch('/api/profile')
  if (res.ok) {
    const data = await res.json()
    Object.assign(form, data)
    Object.assign(preferences, data.preferences || {})
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

  if (!preferences.target_roles || !preferences.location) {
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
    !newEducation.school ||
    !newEducation.degree ||
    !newEducation.field ||
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
      Object.keys(newEducation).forEach(k => newEducation[k] = '')
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

function cancelEdit() { editId.value = null }

async function updateEducation(id) {
  reset('education')

  if (
    !editEducation.school ||
    !editEducation.degree ||
    !editEducation.field ||
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

  if (!editSkill.name) {
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

  if (!newSkill.name) {
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
      Object.keys(newSkill).forEach(k => newSkill[k] = '')
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

function moveSkillUp(index) {
  if (index === 0) return
  const temp = skillsList.value[index]
  skillsList.value[index] = skillsList.value[index - 1]
  skillsList.value[index - 1] = temp
}

function moveSkillDown(index) {
  if (index === skillsList.value.length - 1) return
  const temp = skillsList.value[index]
  skillsList.value[index] = skillsList.value[index + 1]
  skillsList.value[index + 1] = temp
}
</script>

<style scoped src="@/assets/css/profile.css"></style>


