package services

import (
	"dhl/models"
	"dhl/repository"
	"dhl/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

type NotificationService interface {
	SendDistributeAssessmentMail(userIds []string, assessementSeq string, isManager bool) error
	AddUsersToNotify(userIds []*string) error

	RegisterUserInNotify(fcmToken, phone *string, email string) (uuid.UUID, error)
}

type NotificationServiceImpl struct {
	userRepo       repository.UserRepository
	assessmentRepo repository.AssessmentRepository
}

func NewNotificationService(userRepo repository.UserRepository, assessmentRepo repository.AssessmentRepository) NotificationService {
	return &NotificationServiceImpl{userRepo: userRepo, assessmentRepo: assessmentRepo}
}

func (e *NotificationServiceImpl) SendDistributeAssessmentMail(userIds []string, assessementSeq string, isManager bool) error {
	asmt, err := e.assessmentRepo.GetAssessmentMstByAssmtSeq(assessementSeq)
	if err != nil {
		return err
	}

	header := map[string]string{
		"X-API-Key": os.Getenv("NOTIFY_API_KEY"),
	}
	var errorList []error
	for _, uid := range userIds {
		user, _ := e.userRepo.FindByUserId(uid)
		if user.NotifyId == "" {
			log.Println("[Error] Distribution Mail User not registred on notify :", user.UserID, " : ", err)
			errorList = append(errorList, err)
			continue
		}
		sendBody := map[string]interface{}{
			"target_type":   "recipient_id",
			"target_value":  user.NotifyId,
			"template_code": 13,
			"channels":      []string{"email"},
			"data": map[string]interface{}{
				"endDate":        asmt.ValidTo,
				"userName":       fmt.Sprintf("%s %s", user.FirstName, user.LastName),
				"assessmentLink": "https://dhl.catseye.cloud/",
				"isManager":      isManager,
			},
		}
		_, _, sendErr := utils.MakeRESTRequest(http.MethodPost, os.Getenv("NOTIFY_SERVER_URL")+"/api/v1/notifications/send", sendBody, header)
		errorList = append(errorList, sendErr)
	}
	finalErr := errors.Join(errorList...)
	return finalErr
}

func (s *NotificationServiceImpl) RegisterUserInNotify(fcmToken, phone *string, email string) (uuid.UUID, error) {
	header := map[string]string{
		"X-API-Key": os.Getenv("NOTIFY_API_KEY"),
	}
	requestBody := map[string]interface{}{}
	if fcmToken != nil && strings.TrimSpace(*fcmToken) != "" {
		requestBody["fcm_token"] = *fcmToken
	}
	if strings.TrimSpace(email) != "" {
		requestBody["email"] = email
	}
	if phone != nil && strings.TrimSpace(*phone) != "" {
		requestBody["phone"] = *phone
	}
	_, response, sendErr := utils.MakeRESTRequest(http.MethodPost, os.Getenv("NOTIFY_SERVER_URL")+"/api/v1/recipient/register", requestBody, header)
	if sendErr != nil {
		return uuid.Nil, sendErr
	}
	recipient_id, err := utils.ExtractRecipientID(response)
	if err != nil {
		return uuid.Nil, err
	}
	return recipient_id, nil
}

func (e *NotificationServiceImpl) AddUsersToNotify(userIds []*string) error {
	var errorList []error
	users, err := e.userRepo.FetchAllUsers(userIds)
	if err != nil {
		return err
	}
	for _, user := range users {
		recipientId, err := e.RegisterUserInNotify(nil, &user.Phone, user.Email)
		if err != nil {
			log.Println("[Error] RegisterUserInNotify :", user.UserID, " : ", err)
			errorList = append(errorList, err)
			continue
		}
		notifyId := recipientId.String()
		updateData := models.UserProfileUpdate{
			NotifyId: &notifyId,
		}

		err = e.userRepo.UpdateUserProfile(user.UserID, updateData)
		if err != nil {
			log.Println("[Error] userRepo Update User:", user.Email, " : ", err)
			errorList = append(errorList, err)
			continue
		}
	}
	finalErr := errors.Join(errorList...)
	return finalErr
}
