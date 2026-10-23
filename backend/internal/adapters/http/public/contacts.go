package public

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/app/contactimport"
	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
	"github.com/Esca6585dev/habarchy/backend/internal/app/otp"
)

// importTemplate documents the expected upload format and serves a downloadable
// sample. GET ?format=csv or ?format=xlsx returns the file; otherwise JSON
// describing the columns, so the UI can show "what the file should look like"
// before the user picks one.
func (h *Handlers) importTemplate(c *fiber.Ctx) error {
	switch strings.ToLower(c.Query("format")) {
	case "csv":
		c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")
		c.Set(fiber.HeaderContentDisposition, `attachment; filename="habarchy-contacts-template.csv"`)
		return c.Send(contactimport.SampleCSV())
	case "xlsx":
		b, err := contactimport.SampleXLSX()
		if err != nil {
			return err
		}
		c.Set(fiber.HeaderContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Set(fiber.HeaderContentDisposition, `attachment; filename="habarchy-contacts-template.xlsx"`)
		return c.Send(b)
	}
	return httpx.OK(c, fiber.Map{
		"columns":  contactimport.TemplateColumns(),
		"formats":  []string{"xlsx", "csv", "vcf", "docx", "txt"},
		"max_rows": contactimport.MaxRows,
		"sample": fiber.Map{
			"csv":  "/api/v1/contacts/import/template?format=csv",
			"xlsx": "/api/v1/contacts/import/template?format=xlsx",
		},
		"notes": "Birinji setir sütün atlary bolmaly, iň bolmanda `phone`. Türkmen nomerler +993... ýa-da 8 belgili (65123456). Iň köp 5000 setir.",
	})
}

func otpInput(req otpSendRequest, ip string) otp.SendInput {
	return otp.SendInput{Channel: req.Channel, To: req.To, Length: req.Length, TTL: time.Duration(req.TTL) * time.Second,
		Template: req.Template, Locale: req.Locale, Data: req.Data, ClientIP: ip}
}

func contactInput(req contactRequest) contacts.Input {
	return contacts.Input{ExternalID: req.ExternalID, Name: req.Name, Phone: req.Phone, Email: req.Email, WhatsApp: req.WhatsApp,
		TelegramChatID: req.TelegramChatID, SlackID: req.SlackID, Locale: req.Locale, Tags: req.Tags, Attributes: req.Attributes}
}

func (h *Handlers) listContacts(c *fiber.Ctx) error {
	page := httpx.ParsePage(c, 50, 200)
	rows, err := h.Contacts.List(c.UserContext(), caller(c).Project.ID, c.Query("tag"), c.Query("search"), page.Limit, page.Offset)
	if err != nil {
		return err
	}
	out := make([]ContactResponse, 0, len(rows))
	for i := range rows {
		out = append(out, ToContact(&rows[i]))
	}
	return httpx.JSON(c, fiber.StatusOK, out, fiber.Map{"limit": page.Limit, "offset": page.Offset})
}

func (h *Handlers) createContact(c *fiber.Ctx) error {
	var req contactRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	// Clients syncing users send external_id; make the call idempotent.
	ct, created, err := h.Contacts.Upsert(c.UserContext(), caller(c).Project.ID, contactInput(req))
	if err != nil {
		return err
	}
	status := fiber.StatusOK
	if created {
		status = fiber.StatusCreated
	}
	return httpx.JSON(c, status, ToContact(ct), nil)
}

func (h *Handlers) getContact(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "contact_id")
	if err != nil {
		return err
	}
	ct, err := h.Contacts.Get(c.UserContext(), caller(c).Project.ID, id)
	if err != nil {
		return err
	}
	return httpx.OK(c, ToContact(ct))
}

func (h *Handlers) updateContact(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "contact_id")
	if err != nil {
		return err
	}
	var req contactRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	ct, err := h.Contacts.Update(c.UserContext(), caller(c).Project.ID, id, contactInput(req))
	if err != nil {
		return err
	}
	return httpx.OK(c, ToContact(ct))
}

func (h *Handlers) deleteContact(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "contact_id")
	if err != nil {
		return err
	}
	if err := h.Contacts.Delete(c.UserContext(), caller(c).Project.ID, id); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

func (h *Handlers) registerDevice(c *fiber.Ctx) error {
	var req deviceRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	d, err := h.Contacts.RegisterDevice(c.UserContext(), caller(c).Project.ID, contacts.DeviceInput{
		Token: req.Token, Platform: req.Platform, AppVersion: req.AppVersion, ContactID: req.ContactID, ExternalID: req.ExternalID,
	})
	if err != nil {
		return err
	}
	return httpx.Created(c, ToDevice(d))
}

func (h *Handlers) listDevices(c *fiber.Ctx) error {
	page := httpx.ParsePage(c, 50, 200)
	rows, err := h.Contacts.ListDevices(c.UserContext(), caller(c).Project.ID, page.Limit, page.Offset)
	if err != nil {
		return err
	}
	out := make([]DeviceResponse, 0, len(rows))
	for i := range rows {
		out = append(out, ToDevice(&rows[i]))
	}
	return httpx.OK(c, out)
}

func (h *Handlers) removeDevice(c *fiber.Ctx) error {
	if err := h.Contacts.RemoveDevice(c.UserContext(), caller(c).Project.ID, c.Params("token")); err != nil {
		return err
	}
	return httpx.NoContent(c)
}
