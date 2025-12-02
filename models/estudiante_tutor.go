package models

type EstudianteTutor struct {
	EstudianteID uint       `json:"estudiante_id" gorm:"not null"`
	TutorID      uint       `json:"tutor_id" gorm:"not null"`
	Estudiante   Estudiante `json:"estudiante" gorm:"foreignKey:EstudianteID"`
	Tutor        Tutor      `json:"tutor" gorm:"foreignKey:TutorID"`
}
