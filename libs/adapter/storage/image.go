package storage

import (
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
)

const ImageTable string = "images"

// Project Service type for interface in interfaces folder
type ImageStorage struct {
	Session *r.Session
	storage.IImage
}

func (i *ImageStorage) GetByID(user, id string) (*model.Image, error) {

	var err error
	var image = new(model.Image)
	var user_filter = r.Row.Field("user").Eq(id)
	res, err := r.Table(ImageTable).Get(id).Filter(user_filter).Run(i.Session)
	if err != nil {
		return nil, err
	}

	if res.IsNil() {
		return nil, nil
	}

	res.One(image)

	defer res.Close()
	return image, nil
}

func (i *ImageStorage) GetByUser(id string) (*model.ImageList, error) {

	var err error
	var images = new(model.ImageList)

	res, err := r.Table(ImageTable).Get(id).Run(i.Session)
	if err != nil {
		return nil, err
	}

	res.All(images)

	defer res.Close()
	return images, nil
}

func (i *ImageStorage) ListByProject(user, id string) (*model.ImageList, error) {

	var err error
	var images = new(model.ImageList)
	var project_filter = r.Row.Field("project").Field("id").Eq(id)
	var user_filter = r.Row.Field("user").Eq(user)

	res, err := r.Table(ImageTable).Filter(project_filter).Filter(user_filter).Run(i.Session)
	if err != nil {
		return nil, err
	}

	if res.IsNil() {
		return nil, nil
	}

	res.All(images)

	defer res.Close()
	return images, nil
}

func (i *ImageStorage) ListByService(user, id string) (*model.ImageList, error) {

	var err error
	var images = new(model.ImageList)

	var project_filter = r.Row.Field("project").Field("id").Eq(id)
	var user_filter = r.Row.Field("user").Eq(user)
	res, err := r.Table(ImageTable).Filter(project_filter).Filter(user_filter).Run(i.Session)
	if err != nil {
		return nil, err
	}

	res.All(images)

	defer res.Close()
	return images, nil
}

// Insert new image into storage
func (i *ImageStorage) Insert(image *model.Image) (*model.Image, error) {

	res, err := r.Table(ImageTable).Insert(image, r.InsertOpts{ReturnChanges: true}).RunWrite(i.Session)
	if err != nil {
		return nil, err
	}

	image.ID = res.GeneratedKeys[0]

	return image, nil
}

// Update build model
func (i *ImageStorage) Update(image *model.Image) (*model.Image, error) {
	var user_filter = r.Row.Field("user").Eq(image.User)
	_, err := r.Table(ImageTable).Get(image.ID).Filter(user_filter).Replace(image, r.ReplaceOpts{ReturnChanges: true}).RunWrite(i.Session)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func newImageStorage(session *r.Session) *ImageStorage {
	r.TableCreate(ImageTable, r.TableCreateOpts{}).Run(session)
	s := new(ImageStorage)
	s.Session = session
	return s
}
