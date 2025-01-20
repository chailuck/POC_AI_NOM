// internal/handlers/handlers.go (continued)
        return c.JSON(http.StatusInternalServerError, Response{
            Code:    http.StatusInternalServerError,
            Message: "Failed to create individual",
        })
    }

    h.Logger.Infow("Successfully created individual",
        "id", individual.ID,
        "duration", time.Since(start))
    
    return c.JSON(http.StatusCreated, individual)
}

func (h *Handler) GetIndividual(c echo.Context) error {
    start := time.Now()
    id := c.Param("id")
    h.Logger.Infow("Starting GetIndividual request", "id", id)

    var individual models.Individual
    if err := h.DB.Preload("ContactMedium").
        Preload("ExternalReference").
        Preload("IndividualIdentification").
        Preload("PartyCharacteristic").
        First(&individual, "id = ?", id).Error; err != nil {
        
        h.Logger.Errorw("Failed to get individual",
            "id", id,
            "error", err,
            "duration", time.Since(start))
            
        if err == gorm.ErrRecordNotFound {
            return c.JSON(http.StatusNotFound, Response{
                Code:    http.StatusNotFound,
                Message: "Individual not found",
            })
        }
        
        return c.JSON(http.StatusInternalServerError, Response{
            Code:    http.StatusInternalServerError,
            Message: "Failed to get individual",
        })
    }

    h.Logger.Infow("Successfully retrieved individual",
        "id", id,
        "duration", time.Since(start))
        
    return c.JSON(http.StatusOK, individual)
}

func (h *Handler) UpdateIndividual(c echo.Context) error {
    start := time.Now()
    id := c.Param("id")
    h.Logger.Infow("Starting UpdateIndividual request", "id", id)

    // First check if individual exists
    var existingIndividual models.Individual
    if err := h.DB.First(&existingIndividual, "id = ?", id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return c.JSON(http.StatusNotFound, Response{
                Code:    http.StatusNotFound,
                Message: "Individual not found",
            })
        }
        return c.JSON(http.StatusInternalServerError, Response{
            Code:    http.StatusInternalServerError,
            Message: "Failed to check individual existence",
        })
    }

    var updateIndividual models.Individual
    if err := c.Bind(&updateIndividual); err != nil {
        h.Logger.Errorw("Failed to bind request body",
            "error", err,
            "duration", time.Since(start))
        return c.JSON(http.StatusBadRequest, Response{
            Code:    http.StatusBadRequest,
            Message: "Invalid request body",
        })
    }

    updateIndividual.ID = id
    updateIndividual.ModificationDate = time.Now()
    
    // Start a transaction
    tx := h.DB.Begin()

    // Update main individual record
    if err := tx.Model(&existingIndividual).Updates(updateIndividual).Error; err != nil {
        tx.Rollback()
        h.Logger.Errorw("Failed to update individual",
            "id", id,
            "error", err,
            "duration", time.Since(start))
        return c.JSON(http.StatusInternalServerError, Response{
            Code:    http.StatusInternalServerError,
            Message: "Failed to update individual",
        })
    }

    // Update related records if provided
    if len(updateIndividual.ContactMedium) > 0 {
        if err := tx.Model(&existingIndividual).Association("ContactMedium").Replace(updateIndividual.ContactMedium); err != nil {
            tx.Rollback()
            return c.JSON(http.StatusInternalServerError, Response{
                Code:    http.StatusInternalServerError,
                Message: "Failed to update contact medium",
            })
        }
    }

    if len(updateIndividual.ExternalReference) > 0 {
        if err := tx.Model(&existingIndividual).Association("ExternalReference").Replace(updateIndividual.ExternalReference); err != nil {
            tx.Rollback()
            return c.JSON(http.StatusInternalServerError, Response{
                Code:    http.StatusInternalServerError,
                Message: "Failed to update external references",
            })
        }
    }

    if len(updateIndividual.IndividualIdentification) > 0 {
        if err := tx.Model(&existingIndividual).Association("IndividualIdentification").Replace(updateIndividual.IndividualIdentification); err != nil {
            tx.Rollback()
            return c.JSON(http.StatusInternalServerError, Response{
                Code:    http.StatusInternalServerError,
                Message: "Failed to update individual identification",
            })
        }
    }

    if len(updateIndividual.PartyCharacteristic) > 0 {
        if err := tx.Model(&existingIndividual).Association("PartyCharacteristic").Replace(updateIndividual.PartyCharacteristic); err != nil {
            tx.Rollback()
            return c.JSON(http.StatusInternalServerError, Response{
                Code:    http.StatusInternalServerError,
                Message: "Failed to update party characteristics",
            })
        }
    }

    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        return c.JSON(http.StatusInternalServerError, Response{
            Code:    http.StatusInternalServerError,
            Message: "Failed to commit updates",
        })
    }

    h.Logger.Infow("Successfully updated individual",
        "id", id,
        "duration", time.Since(start))
        
    return c.JSON(http.StatusOK, updateIndividual)
}

func (h *Handler) DeleteIndividual(c echo.Context) error {
    start := time.Now()
    id := c.Param("id")
    h.Logger.Infow("Starting DeleteIndividual request", "id", id)

    tx := h.DB.Begin()

    // Delete associated records first
    if err := tx.Where("individual_id = ?", id).Delete(&models.ContactMedium{}).Error; err != nil {
        tx.Rollback()
        return c.JSON(http.StatusInternalServerError, Response{
            Code:    http.StatusInternalServerError,
            Message: "Failed to delete contact medium",
        })
    }

    if err := tx.Where("individual_id = ?", id).Delete(&models.ExternalReference{}).Error; err != nil {
        tx.Rollback()
        return c.JSON(http.StatusInternalServerError, Response{
            Code:    http.StatusInternalServerError,
            Message: "Failed to delete external references",
        })
    }

    if err := tx.Where("individual_id = ?", id).Delete(&models.IndividualIdentification{}).Error; err != nil {
        tx.Rollback()
        return c.JSON(http.StatusInternalServerError, Response{
            Code:    http.StatusInternalServerError,
            Message: "Failed to delete individual identification",
        })
    }

    if err := tx.Where("individual_id = ?", id).Delete(&models.PartyCharacteristic{}).Error; err != nil {
        tx.Rollback()
        return c.JSON(http.StatusInternalServerError, Response{
            Code:    http.StatusInternalServerError,
            Message: "Failed to delete party characteristics",
        })
    }

    // Delete the individual record
    if err := tx.Delete(&models.Individual{}, "id = ?", id).Error; err != nil {
        tx.Rollback()
        h.Logger.Errorw("Failed to delete individual",
            "id", id,
            "error", err,
            "duration", time.Since(start))
        return c.JSON(http.StatusInternalServerError, Response{
            Code:    http.StatusInternalServerError,
            Message: "Failed to delete individual",
        })
    }

    if err := tx.Commit().Error; err != nil {
        return c.JSON(http.StatusInternalServerError, Response{
            Code:    http.StatusInternalServerError,
            Message: "Failed to commit deletion",
        })
    }

    h.Logger.Infow("Successfully deleted individual",
        "id", id,
        "duration", time.Since(start))
        
    return c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) ListIndividuals(c echo.Context) error {
    start := time.Now()
    h.Logger.Info("Starting ListIndividuals request")

    var individuals []models.Individual
    if err := h.DB.Find(&individuals).Error; err != nil {
        h.Logger.Errorw("Failed to list individuals",
            "error", err,
            "duration", time.Since(start))
        return c.JSON(http.StatusInternalServerError, Response{
            Code:    http.StatusInternalServerError,
            Message: "Failed to list individuals",
        })
    }

    h.Logger.Infow("Successfully listed individuals",
        "count", len(individuals),
        "duration", time.Since(start))
        
    return c.JSON(http.StatusOK, individuals)
}
